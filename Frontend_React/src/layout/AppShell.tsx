import { Link, useLocation, useNavigate } from "react-router-dom";

export default function AppShell({ children }:{ children: React.ReactNode }) {
  const { pathname } = useLocation();
  const nav = useNavigate();

  function logout() {
    localStorage.removeItem("token");
    nav("/login", { replace: true });
  }

  return (
    <div className="min-h-screen grid lg:grid-cols-[240px_1fr]">
      {/* Sidebar */}
      <aside className="hidden lg:block bg-white border-r border-gray-100">
        <div className="p-4 text-xl font-semibold">Admin</div>
        <nav className="px-2 space-y-1">
          <Link
            className={`block rounded-xl px-3 py-2 ${pathname==="/admin"?"bg-brand-50 text-brand-700":"hover:bg-gray-100"}`}
            to="/admin"
          >
            Người dùng
          </Link>
        </nav>
      </aside>

      {/* Main */}
      <div className="flex flex-col">
        <header className="bg-white/90 backdrop-blur border-b border-gray-100 h-14 flex items-center justify-between px-4 sticky top-0 z-20">
          <div className="font-semibold">Bảng điều khiển</div>
          <button className="btn-ghost" onClick={logout}>Đăng xuất</button>
        </header>
        <main className="p-4 sm:p-6">{children}</main>
      </div>
    </div>
  );
}
