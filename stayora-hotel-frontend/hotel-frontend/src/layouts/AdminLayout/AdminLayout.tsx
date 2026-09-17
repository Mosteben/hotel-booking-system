import { useState } from "react";
import { Outlet } from "react-router-dom";
import { AdminDataProvider } from "@/context/AdminDataContext";
import { AdminSidebar } from "./AdminSidebar";
import { AdminTopbar } from "./AdminTopbar";

// The Admin area's own shell - deliberately a different structure from the
// customer site's Navbar+Footer (dark sidebar, topbar, dense content),
// not the customer chrome with extra links bolted on.
export function AdminLayout() {
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);

  return (
    <AdminDataProvider>
      <div className="flex min-h-screen bg-bg">
        <AdminSidebar
          mobileOpen={mobileSidebarOpen}
          onClose={() => setMobileSidebarOpen(false)}
        />
        <div className="flex min-w-0 flex-1 flex-col">
          <AdminTopbar onMenuClick={() => setMobileSidebarOpen(true)} />
          <main className="flex-1 px-4 py-6 sm:px-6 lg:px-8">
            <Outlet />
          </main>
        </div>
      </div>
    </AdminDataProvider>
  );
}
