import React, { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { me } from "../api/auth";
import type { Role } from "../types";

export function Protected({
  children,
  allow = ["user", "admin"],
}: {
  children: React.ReactNode;
  allow?: Role[];
}) {
  const [ok, setOk] = useState(false);
  const [checked, setChecked] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        const u = await me();
        setOk(allow.includes(u.role as Role));
      } catch {
        setOk(false);
      } finally {
        setChecked(true);
      }
    })();
  }, [allow]);

  if (!checked) return null;
  if (!ok) return <Navigate to="/login" replace />;
  return <>{children}</>;
}
