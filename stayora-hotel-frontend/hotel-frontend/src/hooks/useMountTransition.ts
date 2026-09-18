import { useEffect, useState } from "react";

// Keeps a component mounted for `exitDuration`ms after `active` goes
// false, so a CSS exit animation (modal-panel-out, modal-backdrop-out)
// actually gets to play instead of the element just vanishing instantly.
// Used by Modal/ConfirmDialog instead of a library - the animation itself
// still lives in index.css as plain keyframes, this just times the unmount.
export function useMountTransition(active: boolean, exitDuration = 180) {
  const [shouldMount, setShouldMount] = useState(active);

  // Entering: reflect immediately during render rather than in an effect -
  // this is React's own recommended pattern for deriving state from a
  // changed prop without an extra render/effect round-trip.
  if (active && !shouldMount) {
    setShouldMount(true);
  }

  // Exiting genuinely needs a real timer, so it belongs in an effect - the
  // setState call here runs inside the timeout callback, not synchronously
  // in the effect body itself.
  useEffect(() => {
    if (active) return;
    const timer = setTimeout(() => setShouldMount(false), exitDuration);
    return () => clearTimeout(timer);
  }, [active, exitDuration]);

  return shouldMount;
}
