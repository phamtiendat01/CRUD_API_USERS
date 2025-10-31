import { useEffect, useState } from "react";
import { me } from "../api/auth";
import { Navigate } from "react-router-dom";

export function Protected({ children, allow=["user","admin"]}:{children:JSX.Element; allow?:("user"|"admin")[]}) {
  const [state,setState]=useState<{loading:boolean; role?:string}>({loading:true});
  useEffect(()=>{ me().then(u=>setState({loading:false, role:u.role})).catch(()=>setState({loading:false})); },[]);
  if (state.loading) return null;
  if (!state.role || !allow.includes(state.role as any)) return <Navigate to="/login" replace />;
  return children;
}
