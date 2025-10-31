import { useEffect, useMemo, useState } from "react";
import { getUsers, createUser, updateUser, deleteUser } from "../../api/admin";
import type { Role, Status, User } from "../../types";
// Nếu file ../../api/admin export các kiểu input, bỏ comment 2 dòng dưới để type chặt chẽ hơn:
// import type { CreateUserInput, UpdateUserInput } from "../../api/admin";
import { fmtDate } from "../../utils/format";
import AppShell from "../../layout/AppShell";
import Avatar from "../../components/ui/Avatar";
import Modal from "../../components/ui/Modal";
import Select from "../../components/ui/Select";
import UserForm from "./UserForm";

/* ===== Icons (inline SVG) ===== */
function IconPlus(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" {...props}>
      <path d="M12 5v14M5 12h14" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
    </svg>
  );
}
function IconReset(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" {...props}>
      <path d="M4 7v6h6M20 17a8 8 0 10-3 3" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    </svg>
  );
}
function IconPencil(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" {...props}>
      <path d="M3 17.25V21h3.75L19.81 7.94a1.5 1.5 0 000-2.12l-1.63-1.63a1.5 1.5 0 00-2.12 0L3 17.25z" stroke="currentColor" strokeWidth="2" strokeLinejoin="round"/>
      <path d="M14.5 5.5l4 4" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
    </svg>
  );
}
function IconTrash(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" {...props}>
      <path d="M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7l1-2h4l1 2" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    </svg>
  );
}

/* ===== Tabs ===== */
const STATUS_TABS: (Status | "all")[] = ["all", "active", "inactive", "banned"];

export default function UsersPage() {
  const [rows, setRows] = useState<User[]>([]);
  const [q, setQ] = useState("");
  const [role, setRole] = useState<Role | "all">("all");
  const [status, setStatus] = useState<Status | "all">("all");
  const [loading, setLoading] = useState(true);

  // modal
  const [openCreate, setOpenCreate] = useState(false);
  const [editUser, setEditUser] = useState<User | null>(null);
  const [delUser, setDelUser] = useState<User | null>(null);

  useEffect(() => {
    setLoading(true);
    getUsers().then(setRows).finally(() => setLoading(false));
  }, []);

  const statusCount = useMemo(() => {
    const m: Record<Status | "all", number> = {
      all: rows.length,
      active: 0,
      inactive: 0,
      banned: 0,
    };
    rows.forEach((r) => {
      m[r.status] = (m[r.status] || 0) + 1;
    });
    return m;
  }, [rows]);

  const data = useMemo(() => {
    const qLower = q.toLowerCase();
    return rows.filter((u) => {
      const matchQ = [u.username, u.full_name, u.email, u.phone]
        .filter(Boolean)
        .some((v) => (v || "").toLowerCase().includes(qLower));
      const matchRole = role === "all" || u.role === role;
      const matchStatus = status === "all" || u.status === status;
      return matchQ && matchRole && matchStatus;
    });
  }, [rows, q, role, status]);

  // Nếu đã import CreateUserInput/UpdateUserInput ở trên, thay `any` bằng các kiểu đó
  async function handleCreate(payload: any /* CreateUserInput */) {
    const u = await createUser(payload);
    setRows((s) => [u, ...s]);
    setOpenCreate(false);
  }
  async function handleUpdate(payload: any /* UpdateUserInput */) {
    if (!editUser) return;
    const u = await updateUser(editUser.id, payload);
    setRows((s) => s.map((x) => (x.id === u.id ? u : x)));
    setEditUser(null);
  }
  async function handleDelete() {
    if (!delUser) return;
    await deleteUser(delUser.id);
    setRows((s) => s.filter((x) => x.id !== delUser.id));
    setDelUser(null);
  }

  function resetFilters() {
    setQ("");
    setRole("all");
    setStatus("all");
  }

  const iconBtn =
    "h-9 w-9 inline-grid place-items-center rounded-lg border text-gray-600 bg-white " +
    "border-gray-200 hover:bg-gray-50 hover:text-gray-900 transition shadow-sm";
  const iconBtnPrimary =
    "h-9 w-9 inline-grid place-items-center rounded-lg border text-white " +
    "bg-indigo-600 border-indigo-600 hover:bg-indigo-700 transition shadow-sm";

  return (
    <AppShell>
      {/* Header */}
      <div className="mb-3">
        <h1 className="text-2xl font-semibold">Quản lý người dùng</h1>
        <p className="text-sm text-gray-500">Theo dõi & thao tác với tài khoản hệ thống.</p>
      </div>

      {/* Tabs trạng thái */}
      <div className="flex flex-wrap gap-2 mb-3">
        {STATUS_TABS.map((s) => {
          const active = status === s;
          return (
            <button
              key={s}
              onClick={() => setStatus(s)}
              className={`px-3 py-1.5 rounded-full border text-sm transition-all
                ${
                  active
                    ? "bg-indigo-50 text-indigo-700 border-indigo-200"
                    : "bg-white text-gray-600 border-gray-200 hover:bg-gray-50 hover:text-gray-800"
                }`}
            >
              <span className="capitalize">{s === "all" ? "Tất cả" : s}</span>
              <span
                className={`ml-2 inline-block px-2 py-0.5 rounded-full text-xs
                ${active ? "bg-indigo-100 text-indigo-700" : "bg-gray-100 text-gray-600"}`}
              >
                {statusCount[s]}
              </span>
            </button>
          );
        })}
      </div>

      {/* Filter + 2 icon bên phải */}
      <div className="card p-3 mb-4">
        <div className="grid grid-cols-1 md:grid-cols-5 gap-3 items-center">
          <input
            className="input md:col-span-2"
            placeholder="Mã / Tên / SĐT / Email"
            value={q}
            onChange={(e) => setQ(e.target.value)}
          />
          <Select value={role} onChange={(e) => setRole(e.target.value as any)}>
            <option value="all">Tất cả vai trò</option>
            <option value="user">user</option>
            <option value="admin">admin</option>
          </Select>
          <Select value={status} onChange={(e) => setStatus(e.target.value as any)}>
            <option value="all">Tất cả trạng thái</option>
            <option value="active">active</option>
            <option value="inactive">inactive</option>
            <option value="banned">banned</option>
          </Select>

          <div className="flex justify-end gap-2">
            <button className={iconBtn} aria-label="Reset" title="Reset" onClick={resetFilters}>
              <IconReset className="h-5 w-5" />
            </button>
            <button
              className={iconBtnPrimary}
              aria-label="Thêm người dùng"
              title="Thêm người dùng"
              onClick={() => setOpenCreate(true)}
            >
              <IconPlus className="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      {/* Bảng */}
      <div className="card overflow-hidden">
        <table className="w-full table-auto text-sm">
          <thead className="text-xs text-gray-600 uppercase bg-gray-50 sticky top-0">
            <tr>
              <th className="px-4 py-3 text-left">ID</th>
              <th className="px-4 py-3 text-left">Khách</th>
              <th className="px-4 py-3 text-left">Email</th>
              <th className="px-4 py-3 text-left">SĐT</th>
              <th className="px-4 py-3 text-left">Giới tính</th>
              <th className="px-4 py-3 text-left">Ngày sinh</th>
              <th className="px-4 py-3 text-left">Trạng thái</th>
              <th className="px-4 py-3 text-left">Hành động</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {loading && (
              <tr>
                <td colSpan={8} className="px-4 py-8 text-center text-gray-500">
                  Đang tải dữ liệu…
                </td>
              </tr>
            )}

            {!loading &&
              data.map((u) => (
                <tr key={u.id} className="hover:bg-indigo-50/40 transition">
                  <td className="px-4 py-3">{u.id}</td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <div className="w-9 h-9 rounded-full overflow-hidden ring-1 ring-gray-200 shadow-sm">
                        <Avatar src={u.avatar_url} name={u.full_name || u.username} />
                      </div>
                      <div>
                        <div className="font-medium">{u.full_name || u.username}</div>
                        <div className="text-xs text-gray-500">@{u.username}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-3 break-all">{u.email}</td>
                  <td className="px-4 py-3">{u.phone}</td>
                  <td className="px-4 py-3 capitalize">{u.gender || ""}</td>
                  <td className="px-4 py-3">{fmtDate(u.date_of_birth)}</td>
                  <td className="px-4 py-3">
                    <span
                      className={`px-2.5 py-1 rounded-full text-xs
                    ${
                      u.status === "active"
                        ? "bg-emerald-50 text-emerald-700 ring-1 ring-emerald-200"
                        : u.status === "inactive"
                        ? "bg-amber-50 text-amber-700 ring-1 ring-amber-200"
                        : "bg-rose-50 text-rose-700 ring-1 ring-rose-200"
                    }`}
                    >
                      {u.status}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <button
                        className={iconBtn}
                        aria-label="Sửa"
                        title="Sửa"
                        onClick={() => setEditUser(u)}
                      >
                        <IconPencil className="h-5 w-5" />
                      </button>
                      <button
                        className={`${iconBtn} text-rose-600 border-rose-200 hover:text-rose-700 hover:bg-rose-50`}
                        aria-label="Xóa"
                        title="Xóa"
                        onClick={() => setDelUser(u)}
                      >
                        <IconTrash className="h-5 w-5" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}

            {!loading && data.length === 0 && (
              <tr>
                <td colSpan={8} className="px-4 py-10 text-center text-gray-500">
                  Không có dữ liệu
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Create */}
      <Modal open={openCreate} onClose={() => setOpenCreate(false)} title="Thêm người dùng" size="xl">
        <UserForm mode="create" onSubmit={handleCreate} onCancel={() => setOpenCreate(false)} />
      </Modal>

      {/* Edit */}
      <Modal open={!!editUser} onClose={() => setEditUser(null)} title="Chỉnh sửa người dùng" size="xl">
        {editUser && (
          <UserForm mode="edit" initial={editUser} onSubmit={handleUpdate} onCancel={() => setEditUser(null)} />
        )}
      </Modal>

      {/* Delete */}
      <Modal open={!!delUser} onClose={() => setDelUser(null)} title="Xóa người dùng">
        <p>
          Bạn chắc chắn muốn xóa <b>{delUser?.username}</b>?
        </p>
        <div className="mt-5 flex justify-end gap-2">
          <button
            className="h-9 px-3 rounded-lg bg-white text-gray-700 border border-gray-200 hover:bg-gray-50"
            onClick={() => setDelUser(null)}
          >
            Hủy
          </button>
          <button className="h-9 px-3 rounded-lg bg-rose-600 text-white hover:bg-rose-700" onClick={handleDelete}>
            Xóa
          </button>
        </div>
      </Modal>
    </AppShell>
  );
}
