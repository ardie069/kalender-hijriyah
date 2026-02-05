import { useState, useEffect } from "react";

export function useMounted() {
  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    requestAnimationFrame(() => setMounted(true));
  }, []);
  return mounted;
}
