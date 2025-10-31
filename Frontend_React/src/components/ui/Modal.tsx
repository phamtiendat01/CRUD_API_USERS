export default function Modal({
  open, onClose, title, children, footer, size = "md",
}: {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  footer?: React.ReactNode;
  size?: "sm" | "md" | "lg" | "xl";
}) {
  if (!open) return null;
  const widths = { sm: "max-w-md", md: "max-w-2xl", lg: "max-w-3xl", xl: "max-w-5xl" } as const;

  return (
    <div className="fixed inset-0 z-50">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <div className={`absolute inset-0 grid place-items-center p-4`}>
        <div className={`w-full ${widths[size]} bg-white rounded-2xl shadow-xl border border-gray-100`}>
          <div className="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
            <h3 className="font-semibold">{title}</h3>
            <button className="btn-ghost" onClick={onClose}>✕</button>
          </div>
          <div className="p-5">{children}</div>
          {footer && <div className="px-5 py-4 border-t border-gray-100 bg-gray-50 rounded-b-2xl">{footer}</div>}
        </div>
      </div>
    </div>
  );
}
