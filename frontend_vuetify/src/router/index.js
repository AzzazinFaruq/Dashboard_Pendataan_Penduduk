/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 */

// Composables
import { createRouter, createWebHistory } from "vue-router/auto";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
});

// Halaman yang boleh diakses tanpa login.
const publicPaths = ["/", "/login", "/register", "/forbidden"];

// Guard global: cegah akses halaman terproteksi tanpa auth,
// dan cegah user yang sudah login membuka halaman login/register lagi.
router.beforeEach((to) => {
  const isAuthenticated = localStorage.getItem("auth") === "true";
  const isPublic = publicPaths.includes(to.path);

  if (!isAuthenticated && !isPublic) {
    return { path: "/login" };
  }

  if (isAuthenticated && (to.path === "/login" || to.path === "/register")) {
    return { path: "/dashboard" };
  }

  return true;
});

export default router;
