import { useEffect, useMemo, useState } from "react";
import Input from "../../components/ui/Input";
import Select from "../../components/ui/Select";
import Button from "../../components/ui/Button";
import type { User } from "../../types";
import type { CreateUserInput, UpdateUserInput } from "../../api/admin";

type Mode = "create" | "edit";

function toInputDate(iso?: string) {
  if (!iso) return "";
  const d = new Date(iso);
  return Number.isNaN(d.getTime()) ? "" : d.toISOString().slice(0, 10);
}

export default function UserForm({
  mode, initial, onSubmit, onCancel,
}: {
  mode: Mode;
  initial?: User;
  onSubmit: (payload: CreateUserInput | UpdateUserInput) => Promise<void> | void;
  onCancel: () => void;
}) {
  const [v, setV] = useState<{
    username: string; email: string; password: string; confirm: string;
    full_name?: string; phone?: string; gender?: string; date_of_birth?: string;
    avatar_url?: string; street?: string; city?: string; state?: string;
    country?: string; postal_code?: string; role?: "user"|"admin"; status?: string;
  }>({
    username: "", email: "", password: "", confirm: "",
    full_name: "", phone: "", gender: "", date_of_birth: "",
    avatar_url: "", street: "", city: "", state: "",
    country: "", postal_code: "", role: "user", status: "active",
  });

  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [err, setErr] = useState("");
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    if (mode === "edit" && initial) {
      setV({
        username: initial.username || "", email: initial.email || "",
        password: "", confirm: "",
        full_name: initial.full_name || "", phone: initial.phone || "",
        gender: initial.gender || "", date_of_birth: toInputDate(initial.date_of_birth),
        avatar_url: initial.avatar_url || "", street: initial.street || "",
        city: initial.city || "", state: initial.state || "",
        country: initial.country || "", postal_code: initial.postal_code || "",
        role: (initial.role as any) || "user", status: (initial.status as any) || "active",
      });
      setPreview(initial.avatar_url || null);
    }
  }, [mode, initial]);

  function validate(): string | null {
    if (v.username.length < 3) return "Username tối thiểu 3 ký tự.";
    if (!/^\S+@\S+\.\S+$/.test(v.email)) return "Email không hợp lệ.";
    if (mode === "create" && v.password.length < 6) return "Mật khẩu tối thiểu 6 ký tự.";
    if (mode === "create" && v.password !== v.confirm) return "Nhập lại mật khẩu không khớp.";
    return null;
  }

  const payload: CreateUserInput | UpdateUserInput = useMemo(() => {
    const base = {
      username: v.username.trim(), email: v.email.trim(),
      full_name: v.full_name?.trim() || undefined,
      phone: v.phone?.trim() || undefined,
      gender: v.gender || undefined,
      date_of_birth: v.date_of_birth || undefined,
      // avatar_url: nếu bạn triển khai upload file ở BE thì set sau khi upload
      avatar_url: !avatarFile && v.avatar_url ? v.avatar_url : undefined,
      street: v.street?.trim() || undefined, city: v.city?.trim() || undefined,
      state: v.state?.trim() || undefined, country: v.country?.trim() || undefined,
      postal_code: v.postal_code?.trim() || undefined, role: v.role, status: v.status,
    } as UpdateUserInput;

    if (mode === "create") (base as any).password = v.password;
    else if (v.password) (base as any).password = v.password;
    return base;
  }, [v, mode, avatarFile]);

  async function fileToPreview(file: File): Promise<string> {
    return new Promise((resolve) => {
      const img = new Image(); const reader = new FileReader();
      reader.onload = () => { img.src = reader.result as string; };
      img.onload = () => {
        const canvas = document.createElement("canvas");
        const max = 160;
        const ratio = Math.min(max / img.width, max / img.height, 1);
        canvas.width = Math.round(img.width * ratio);
        canvas.height = Math.round(img.height * ratio);
        const ctx = canvas.getContext("2d")!;
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
        resolve(canvas.toDataURL("image/webp", 0.85));
      };
      reader.readAsDataURL(file);
    });
  }

  async function onPickAvatar(f?: File) {
    if (!f) { setAvatarFile(null); setPreview(v.avatar_url || null); return; }
    setAvatarFile(f);
    const thumb = await fileToPreview(f);
    setPreview(thumb);
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErr("");
    const msg = validate();
    if (msg) { setErr(msg); return; }

    // NOTE: BE hiện tại chưa có endpoint upload file.
    // Nếu cần lưu file thực, upload f trước -> nhận URL -> gán payload.avatar_url = URL -> onSubmit(payload).
    try { setBusy(true); await onSubmit(payload); }
    catch { setErr("Không thể lưu. Vui lòng thử lại."); }
    finally { setBusy(false); }
  }

  return (
    <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-12 gap-4">
      {/* trái */}
      <div className="md:col-span-8 space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Username <span className="text-rose-500">*</span></label>
            <Input value={v.username} onChange={e=>setV(s=>({...s, username:e.target.value}))} required minLength={3}/>
          </div>
          <div>
            <label className="label">Email <span className="text-rose-500">*</span></label>
            <Input type="email" value={v.email} onChange={e=>setV(s=>({...s, email:e.target.value}))} required/>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Họ tên</label>
            <Input value={v.full_name} onChange={e=>setV(s=>({...s, full_name:e.target.value}))}/>
          </div>
          <div>
            <label className="label">SĐT</label>
            <Input value={v.phone} onChange={e=>setV(s=>({...s, phone:e.target.value}))}/>
          </div>
          <div>
            <label className="label">Giới tính</label>
            <Select value={v.gender} onChange={e=>setV(s=>({...s, gender:e.target.value}))}>
              <option value="">—</option><option value="male">Male</option>
              <option value="female">Female</option><option value="other">Other</option>
            </Select>
          </div>
          <div>
            <label className="label">Ngày sinh</label>
            <Input type="date" value={v.date_of_birth} onChange={e=>setV(s=>({...s, date_of_birth:e.target.value}))}/>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Đường</label>
            <Input value={v.street} onChange={e=>setV(s=>({...s, street:e.target.value}))}/>
          </div>
          <div>
            <label className="label">Thành phố</label>
            <Input value={v.city} onChange={e=>setV(s=>({...s, city:e.target.value}))}/>
          </div>
          <div>
            <label className="label">Tỉnh/Bang</label>
            <Input value={v.state} onChange={e=>setV(s=>({...s, state:e.target.value}))}/>
          </div>
          <div>
            <label className="label">Quốc gia</label>
            <Input value={v.country} onChange={e=>setV(s=>({...s, country:e.target.value}))}/>
          </div>
          <div>
            <label className="label">Mã bưu chính</label>
            <Input value={v.postal_code} onChange={e=>setV(s=>({...s, postal_code:e.target.value}))}/>
          </div>
          <div>
            <label className="label">Trạng thái</label>
            <Select value={v.status} onChange={e=>setV(s=>({...s, status:e.target.value}))}>
              <option value="active">active</option>
              <option value="inactive">inactive</option>
              <option value="banned">banned</option>
            </Select>
          </div>
          <div>
            <label className="label">Role</label>
            <Select value={v.role} onChange={e=>setV(s=>({...s, role: e.target.value as "user"|"admin"}))}>
              <option value="user">user</option>
              <option value="admin">admin</option>
            </Select>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="label">{mode==="create" ? "Mật khẩu *" : "Mật khẩu (đổi nếu nhập)"}</label>
            <Input type="password" value={v.password} onChange={e=>setV(s=>({...s, password:e.target.value}))} minLength={mode==="create" ? 6 : 0}/>
          </div>
          <div>
            <label className="label">{mode==="create" ? "Nhập lại mật khẩu *" : "Nhập lại mật khẩu"}</label>
            <Input type="password" value={v.confirm} onChange={e=>setV(s=>({...s, confirm:e.target.value}))}/>
          </div>
        </div>

        {err && <p className="text-sm text-rose-600">{err}</p>}

        <div className="flex items-center justify-end gap-2 pt-2">
          <button
            type="submit"
            className="px-3 py-2 rounded-lg bg-white text-gray-700 border border-gray-200 hover:bg-gray-50"
            onClick={onCancel}
          >
            Hủy
          </button>
          <Button type="submit" disabled={busy} className="bg-indigo-600 hover:bg-indigo-700 text-white">
            {mode === "create" ? "Tạo mới" : "Lưu thay đổi"}
          </Button>
        </div>
      </div>

      {/* phải: chọn ảnh + preview */}
      <div className="md:col-span-4">
        <div className="card p-4 sticky top-4">
          <p className="label mb-2">Ảnh đại diện</p>

          <label
            htmlFor="avatar-input"
            className="block w-full p-4 text-center rounded-xl border-2 border-dashed
                       border-gray-200 hover:border-indigo-300 hover:bg-indigo-50/30 cursor-pointer transition"
          >
            <div className="mx-auto w-24 h-24 rounded-full overflow-hidden ring-1 ring-gray-200 shadow-sm mb-3">
              {preview
                ? <img src={preview} className="w-full h-full object-cover" />
                : <div className="w-full h-full grid place-items-center text-gray-400 text-sm">No image</div>}
            </div>
            <div className="text-sm text-gray-600">
              Kéo thả hoặc <span className="text-indigo-600 font-medium">chọn ảnh</span>
            </div>
            <div className="text-xs text-gray-400">PNG/JPG/WebP • Tối đa ~3MB</div>
          </label>
          <input
            id="avatar-input" type="file" accept="image/*" hidden
            onChange={e=>onPickAvatar(e.target.files?.[0] || undefined)}
          />

          <div className="mt-4">
            <p className="label">Hoặc dán URL ảnh</p>
            <Input
              placeholder="https://…"
              value={v.avatar_url}
              onChange={e=>{ setV(s=>({...s, avatar_url:e.target.value})); setPreview(e.target.value || null); }}
            />
          </div>
        </div>
      </div>
    </form>
  );
}
