export default function ChipToggle({
  active, children, onClick, className = "",
}: {active:boolean; children:React.ReactNode; onClick:()=>void; className?:string}) {
  return (
    <button
      onClick={onClick}
      className={`px-3 py-1.5 rounded-full border text-sm hover-lift transition-colors
        ${active ? "bg-indigo-50 text-indigo-700 border-indigo-200"
                 : "bg-white text-gray-600 border-gray-200 hover:bg-gray-50"} ${className}`}
    >
      {children}
    </button>
  );
}
