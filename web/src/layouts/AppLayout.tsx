import { Outlet } from "react-router-dom";
import Header from "../components/Header";
import Sidebar from "../components/Sidebar";

export default function AppLayout() {
  return (
    <>
      <Header />
      <div className="container page-content">
        <Sidebar />
        <Outlet />
      </div>
    </>
  );
}
