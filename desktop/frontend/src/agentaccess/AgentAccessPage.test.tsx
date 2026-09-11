import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { AgentAccessPage } from "./AgentAccessPage";
import type { AgentAccessBridge, AgentAccessWorkspace } from "./model";

const off: AgentAccessWorkspace = {
  stdioDescription: "Available on demand. Stdio opens no network listener.",
  http: { enabled: false, allowedOrigins: [], requestCount: 0 },
  authorities: [
    { name: "Observe", status: "available", description: "Read bounded data." },
    { name: "Operate", status: "future", description: "Unavailable." },
    { name: "Manage", status: "future", description: "Unavailable." },
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
    expect(screen.getAllByText(/Future, unavailable/i)).toHaveLength(2);
    expect(screen.getByText("Configure localhost HTTP").closest("details")).not.toHaveAttribute("open");
    fireEvent.click(screen.getByText("Configure localhost HTTP"));
    fireEvent.click(screen.getByRole("button", { name: /Enable and copy/i }));
    await waitFor(() => expect(api.enable).toHaveBeenCalled());
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
      screen.getByRole("button", { name: /Revoke access/i }),
    ).toBeInTheDocument();
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
