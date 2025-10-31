import { useState } from "react";
import { login } from "../../api/auth";
import { useNavigate, Link } from "react-router-dom";
import Button from "../../components/ui/Button";
import Input from "../../components/ui/Input";

export default function Login() {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [showPw, setShowPw] = useState(false);
  const [err, setErr] = useState("");
  const nav = useNavigate();

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setErr("");
    try {
      const u = await login(identifier, password);
      nav(u.role === "admin" ? "/admin" : "/app", { replace: true });
    } catch {
      setErr("Sai tài khoản hoặc mật khẩu");
    }
  }

  return (
    <div className="relative auth-bg">
      <span className="auth-overlay" />
      <div className="relative z-10 min-h-screen flex items-center justify-center p-4">
        <div className="auth-card fade-in-up">
          {/* Header */}
          <div className="mb-6 text-center">
            <div className="mx-auto mb-2 h-10 w-10 rounded-xl bg-indigo-600 text-white grid place-items-center shadow">
              {/* logo nhỏ */}
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87L18.18 22 12 18.77 5.82 22 7 14.14l-5-4.87 6.91-1.01L12 2z" fill="currentColor"/>
              </svg>
            </div>
            <h1 className="text-2xl font-semibold">Đăng nhập</h1>
            <p className="text-sm text-gray-500">Chào mừng quay lại 👋</p>
          </div>

          {/* Form */}
          <form onSubmit={onSubmit} className="space-y-4">
            <div>
              <label className="label">Email / Username</label>
              <Input
                value={identifier}
                onChange={(e) => setIdentifier(e.target.value)}
                placeholder="you@example.com"
                autoComplete="username"
                required
              />
            </div>

            <div>
              <label className="label">Mật khẩu</label>
              <div className="relative">
                <Input
                  type={showPw ? "text" : "password"}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  autoComplete="current-password"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPw((s) => !s)}
                  className="absolute inset-y-0 right-2.5 my-auto px-2 text-gray-500 hover:text-gray-700"
                  aria-label={showPw ? "Ẩn mật khẩu" : "Hiện mật khẩu"}
                >
                  {showPw ? (
                    /* eye-off */
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                      <path d="M17.94 17.94A10.94 10.94 0 0 1 12 20C7 20 2.73 16.11 1 12c.74-1.63 1.83-3.09 3.17-4.31M9.88 9.88A3 3 0 0 0 12 15a3 3 0 0 0 2.12-5.12M3 3l18 18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    </svg>
                  ) : (
                    /* eye */
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                      <circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth="2"/>
                    </svg>
                  )}
                </button>
              </div>
            </div>

            {err && <p className="text-red-600 text-sm">{err}</p>}

            <Button type="submit" className="w-full">Đăng nhập</Button>

            <p className="text-sm text-gray-500 text-center">
              Chưa có tài khoản?{" "}
              <Link className="text-indigo-600 hover:underline font-medium" to="/register">
                Đăng ký
              </Link>
            </p>
          </form>
        </div>
      </div>
    </div>
  );
}
