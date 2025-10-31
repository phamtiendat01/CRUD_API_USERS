export const fmtDate = (iso?: string) => {
  if (!iso) return "";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  const y = d.getFullYear();
  const m = `${d.getMonth() + 1}`.padStart(2, "0");
  const day = `${d.getDate()}`.padStart(2, "0");
  return `${day}/${m}/${y}`;
};

export const fullAddress = (u: {
  street?: string; city?: string; state?: string; country?: string; postal_code?: string;
}) => {
  return [u.street, u.city, u.state, u.country, u.postal_code].filter(Boolean).join(", ");
};
