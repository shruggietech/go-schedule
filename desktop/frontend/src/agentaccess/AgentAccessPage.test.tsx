import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { AgentAccessPage } from "./AgentAccessPage";
import type { AgentAccessBridge, AgentAccessWorkspace } from "./model";

const off: AgentAccessWorkspace = {
  stdioDescription: "Available on demand. Stdio opens no network listener.",
  http: { enabled: false, allowedOrigins: [], requestCount: 0 },
  authorities: [
    { name: "Observe", status: "available", description: "Read bounded data." },
    { name: "Operate", status: "available", description: "Run, enable, and disable existing tasks." },
    { name: "Manage", status: "available", description: "Create, update, and delete bounded automation definitions." },
  ],
};
const bridge = (workspace = off): AgentAccessBridge => ({
  workspace: vi
    .fn()
    .mockResolvedValue({
      action: "load_agent_access",
      outcome: "accepted",
      message: "",
      workspace,
    }),
  enable: vi
    .fn()
    .mockResolvedValue({
      action: "enable_agent_access",
      outcome: "accepted",
      message: "copied",
      workspace,
    }),
  rotate: vi.fn(),
  revoke: vi.fn(),
  createGrant: vi.fn().mockResolvedValue({ action: "create_agent_grant", outcome: "accepted", message: "Enrollment copied.", workspace }),
  editGrant: vi.fn().mockResolvedValue({ action: "edit_agent_grant", outcome: "accepted", message: "Grant narrowed.", workspace }),
  revokeGrant: vi.fn().mockResolvedValue({ action: "revoke_agent_grant", outcome: "accepted", message: "Grant revoked.", workspace }),
  actions: vi.fn().mockResolvedValue({ action: "load_agent_actions", outcome: "accepted", message: "", actions: [] }),
  openGuide: vi.fn(),
});

describe("Agent Access page", () => {
  it("explains stdio, authority, and enables a named client without a credential field", async () => {
    const api = bridge();
    const { container } = render(
      <AgentAccessPage bridge={api} available refreshToken={1} />,
    );
    expect(
      await screen.findByRole("heading", { name: "Agent Access" }),
    ).toBeInTheDocument();
    expect(screen.getByText(/opens no network listener/i)).toBeInTheDocument();
    expect(screen.queryByText(/Future, unavailable/i)).not.toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Manage" })).toBeInTheDocument();
    expect(screen.getByText("Configure localhost HTTP").closest("details")).not.toHaveAttribute("open");
    fireEvent.click(screen.getByText("Configure localhost HTTP"));
    fireEvent.change(screen.getByLabelText("Permission"), { target: { value: "operate" } });
    fireEvent.click(screen.getByRole("button", { name: /Enable and copy/i }));
    await waitFor(() => expect(api.enable).toHaveBeenCalled());
    expect(api.enable).toHaveBeenCalledWith(expect.objectContaining({ permission: "operate" }));
    expect(container.textContent).not.toContain("top-secret");
    expect(container.querySelector('input[name="credential"]')).toBeNull();
  });
  it("shows active evidence and lifecycle controls", async () => {
    const active = {
      ...off,
      http: {
        enabled: true,
        endpoint: "http://127.0.0.1:43123/mcp",
        allowedOrigins: [],
        credentialFingerprint: "abc",
        enabledAt: "2026-09-09T12:00:00Z",
        clientName: "Codex",
        permission: "operate",
        lastAccessedAt: "2026-09-09T12:01:00Z",
        requestCount: 2,
      },
    };
    render(
      <AgentAccessPage bridge={bridge(active)} available refreshToken={1} />,
    );
    expect(await screen.findByText("Codex")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Rotate credential/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Revoke localhost access/i }),
    ).toBeInTheDocument();
  });
  it("creates a bounded remote grant without rendering its enrollment secret", async () => {
    const api = bridge({ ...off, daemon: { id: "daemon-local", name: "This computer" }, mcpState: "off", transports: [], grants: [] });
    const { container } = render(<AgentAccessPage bridge={api} available refreshToken={1} />);
    fireEvent.click(await screen.findByRole("button", { name: "Grant remote access" }));
    const dialog = screen.getByRole("dialog", { name: "Grant remote MCP access" });
    fireEvent.change(within(dialog).getByLabelText("Client name"), { target: { value: "Build agent" } });
    fireEvent.change(within(dialog).getByLabelText("Authority"), { target: { value: "operate" } });
    fireEvent.change(within(dialog).getByLabelText("Access duration"), { target: { value: "7d" } });
    fireEvent.click(within(dialog).getByRole("button", { name: "Create and copy enrollment" }));
    await waitFor(() => expect(api.createGrant).toHaveBeenCalledWith({ clientName: "Build agent", capability: "operate", duration: "7d" }));
    expect(container.textContent).not.toContain("one-time-secret");
  });
  it("requires an explicit acknowledgement for a non-expiring grant", async () => {
    render(<AgentAccessPage bridge={bridge()} available refreshToken={1} />);
    fireEvent.click(await screen.findByRole("button", { name: "Grant remote access" }));
    const dialog = screen.getByRole("dialog", { name: "Grant remote MCP access" });
    fireEvent.change(within(dialog).getByLabelText("Access duration"), { target: { value: "non-expiring" } });
    expect(within(dialog).getByRole("button", { name: "Create and copy enrollment" })).toBeDisabled();
    fireEvent.change(within(dialog).getByLabelText("Client name"), { target: { value: "Build agent" } });
    fireEvent.click(within(dialog).getByLabelText(/remain active until revoked/i));
    expect(within(dialog).getByRole("button", { name: "Create and copy enrollment" })).toBeEnabled();
  });
  it("offers Manage with optional per-call confirmation", async () => {
    const api = bridge();
    render(<AgentAccessPage bridge={api} available refreshToken={1} />);
    await screen.findByRole("heading", { name: "Agent Access" });
    fireEvent.click(screen.getByText("Configure localhost HTTP"));
    fireEvent.change(screen.getByLabelText("Permission"), { target: { value: "manage" } });
    fireEvent.click(screen.getByLabelText("Require confirmed calls"));
    fireEvent.click(screen.getByRole("button", { name: /Enable and copy/i }));
    await waitFor(() => expect(api.enable).toHaveBeenCalledWith(expect.objectContaining({ permission: "manage", requireConfirmation: true })));
  });
  it("shows an actionable initial load failure instead of indefinite loading", async () => {
    const api = bridge();
    api.workspace = vi
      .fn()
      .mockResolvedValue({
        action: "load_agent_access",
        outcome: "unavailable",
        message: "Check the local scheduler connection.",
      });
    render(<AgentAccessPage bridge={api} available refreshToken={1} />);
    expect(
      await screen.findByRole("heading", { name: "Agent Access unavailable" }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "Try again" }),
    ).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Try again" }));
    await waitFor(() => expect(api.workspace).toHaveBeenCalledTimes(2));
  });
});
