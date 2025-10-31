import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Protected } from "./components/Protected";
import Login from "./pages/auth/Login";
import Register from "./pages/auth/Register";
import UsersPage from "./pages/admin/UsersPage";
import Home from "./pages/user/Home";

export default function App(){
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login/>}/>
        <Route path="/register" element={<Register/>}/>
        <Route path="/admin" element={<Protected allow={["admin"]}><UsersPage/></Protected>}/>
        <Route path="/app" element={<Protected><Home/></Protected>}/>
        <Route path="*" element={<Login/>}/>
      </Routes>
    </BrowserRouter>
  );
}
