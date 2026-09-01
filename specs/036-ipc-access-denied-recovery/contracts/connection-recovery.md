# Connection Recovery Contract

1. Client transport errors expose one of four stable kinds and unwrap to the
   OS or context cause. API status errors retain their existing type.
2. Access-denied text says the daemon rejected the current session and never
   asks only whether the daemon is running.
3. The GUI exposes at most one in-frame connection panel. It always contains
   Retry and Exit and never blocks tab visibility.
4. Repeated model, calendar, and stream failures update that panel. They do not
   create dialogs.
5. A successful coordinated refresh clears the panel. A failed retry updates it.
6. Windows stale-session guidance appears only from verified group/member/token
   evidence. Unknown evidence remains qualified.
7. Background delays follow 2, 4, 8, 16, 30 seconds and remain at 30 seconds;
   context cancellation and Retry interrupt a wait.
8. Named-pipe ACL construction and authorized identities do not change.
