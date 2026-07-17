// Shared client-side logic for PT Ecommindo Jaya Persada website.
// Talks to the Go backend REST API served from the same origin.

const API_BASE = "";

const LS_AUTH = "ecommindo_auth";
const SS_PENDING_ADD = "ecommindo_pending_add";

/* ---------------- API helper ---------------- */

async function apiFetch(path, options = {}) {
  const auth = getAuth();
  const headers = Object.assign(
    { "Content-Type": "application/json" },
    options.headers || {}
  );
  if (auth && auth.token) {
    headers["Authorization"] = "Bearer " + auth.token;
  }

  const res = await fetch(API_BASE + path, Object.assign({}, options, { headers }));
  const isJson = (res.headers.get("content-type") || "").includes("application/json");
  const data = isJson ? await res.json() : null;

  if (!res.ok) {
    if (res.status === 401) clearAuth();
    throw new Error((data && data.error) || "Terjadi kesalahan pada server.");
  }

  return data;
}

/* ---------------- Auth ---------------- */

function getAuth() {
  return JSON.parse(localStorage.getItem(LS_AUTH) || "null");
}

function setAuth(data) {
  localStorage.setItem(LS_AUTH, JSON.stringify(data));
}

function clearAuth() {
  localStorage.removeItem(LS_AUTH);
}

function isLoggedIn() {
  return !!getAuth();
}

function getCurrentUser() {
  const auth = getAuth();
  return auth ? auth.user : null;
}

async function registerUser({ fullName, phone, email, password }) {
  const data = await apiFetch("/api/auth/register", {
    method: "POST",
    body: JSON.stringify({ full_name: fullName, phone, email, password }),
  });
  setAuth(data);
  return data.user;
}

async function loginUser(email, password) {
  const data = await apiFetch("/api/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
  setAuth(data);
  return data.user;
}

function logoutUser() {
  clearAuth();
}

/* ---------------- Services ---------------- */

async function getServices() {
  return apiFetch("/api/services");
}

const SERVICE_IMAGES = {
  "cyber-security": "images/service-security.jpg",
  "server-build": "images/service-server.jpg",
  "software-dev": "images/service-software.jpg",
  "app-dev": "images/service_application.jpg",
  "infra-management": "images/service_manajemen.jpg",
};

function serviceImageTag(service) {
  const src = SERVICE_IMAGES[service.id];
  if (!src) return "";
  return `<div class="service-photo"><img src="${src}" alt="${escapeHtml(service.name)}" loading="lazy" /></div>`;
}

/* ---------------- Cart ---------------- */

async function getCartSummary() {
  if (!isLoggedIn()) return { items: [], total: 0 };
  return apiFetch("/api/cart");
}

async function addToCart(serviceId) {
  return apiFetch("/api/cart", {
    method: "POST",
    body: JSON.stringify({ service_id: serviceId }),
  });
}

async function removeFromCart(serviceId) {
  return apiFetch("/api/cart/" + encodeURIComponent(serviceId), { method: "DELETE" });
}

/* ---------------- Checkout ---------------- */

async function checkoutOrder() {
  return apiFetch("/api/checkout", { method: "POST" });
}

/* ---------------- Navbar / shared UI ---------------- */

async function initNavbar() {
  const authSlot = document.querySelector("[data-auth-slot]");
  const user = getCurrentUser();

  if (authSlot) {
    if (user) {
      authSlot.innerHTML =
        '<span class="user-chip">👤 ' +
        escapeHtml(user.full_name.split(" ")[0]) +
        '</span><button class="btn btn-outline btn-sm" id="logoutBtn" type="button">Logout</button>';
      const logoutBtn = document.getElementById("logoutBtn");
      logoutBtn.addEventListener("click", () => {
        logoutUser();
        showToast("Anda berhasil logout.");
        setTimeout(() => (window.location.href = "index.html"), 600);
      });
    } else {
      authSlot.innerHTML = '<a href="login.html" class="btn btn-primary btn-sm">Login</a>';
    }
  }

  const cartCountEls = document.querySelectorAll("[data-cart-count]");
  if (cartCountEls.length) {
    try {
      const cart = await getCartSummary();
      const count = cart.items.length;
      cartCountEls.forEach((el) => {
        el.textContent = count;
        el.style.display = count > 0 ? "flex" : "none";
      });
    } catch (e) {
      // ignore, badge stays hidden
    }
  }

  const toggle = document.querySelector(".nav-toggle");
  const links = document.querySelector(".nav-links");
  if (toggle && links) {
    toggle.addEventListener("click", () => links.classList.toggle("open"));
  }

  const current = document.body.getAttribute("data-page");
  if (current) {
    document.querySelectorAll(".nav-links a[data-page]").forEach((a) => {
      if (a.getAttribute("data-page") === current) a.classList.add("active");
    });
  }
}

function escapeHtml(str) {
  const div = document.createElement("div");
  div.textContent = str;
  return div.innerHTML;
}

/* ---------------- Toast ---------------- */

let toastTimer = null;
function showToast(message) {
  let toast = document.querySelector(".toast");
  if (!toast) {
    toast = document.createElement("div");
    toast.className = "toast";
    document.body.appendChild(toast);
  }
  toast.textContent = message;
  toast.classList.add("show");
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => toast.classList.remove("show"), 3000);
}

/* ---------------- Guards ---------------- */

function requireLogin(redirectTo) {
  if (!isLoggedIn()) {
    window.location.href = "login.html?redirect=" + encodeURIComponent(redirectTo);
    return false;
  }
  return true;
}

/* ---------------- Motion: scroll reveal ---------------- */

function initReveal() {
  const els = document.querySelectorAll(".reveal:not(.in-view)");
  if (!els.length) return;

  const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  if (reduceMotion) {
    els.forEach((el) => el.classList.add("in-view"));
    return;
  }

  const io = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add("in-view");
          io.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.15 }
  );

  els.forEach((el) => {
    const siblingIndex = Array.prototype.indexOf.call(el.parentElement.children, el);
    el.style.transitionDelay = Math.min(siblingIndex, 5) * 80 + "ms";
    io.observe(el);
  });
}

/* ---------------- Motion: animated counters ---------------- */

function animateCounter(el) {
  const to = parseFloat(el.dataset.countTo);
  const decimals = parseInt(el.dataset.decimals || "0", 10);
  const prefix = el.dataset.prefix || "";
  const suffix = el.dataset.suffix || "";
  const duration = 1400;
  const start = performance.now();

  function tick(now) {
    const progress = Math.min((now - start) / duration, 1);
    const eased = 1 - Math.pow(1 - progress, 3);
    el.textContent = prefix + (to * eased).toFixed(decimals) + suffix;
    if (progress < 1) requestAnimationFrame(tick);
    else el.textContent = prefix + to.toFixed(decimals) + suffix;
  }

  requestAnimationFrame(tick);
}

function initCounters() {
  const counters = document.querySelectorAll("[data-count-to]");
  if (!counters.length) return;

  if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
    counters.forEach((el) => {
      const prefix = el.dataset.prefix || "";
      const suffix = el.dataset.suffix || "";
      const decimals = parseInt(el.dataset.decimals || "0", 10);
      el.textContent = prefix + parseFloat(el.dataset.countTo).toFixed(decimals) + suffix;
    });
    return;
  }

  const io = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          animateCounter(entry.target);
          io.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.4 }
  );

  counters.forEach((el) => io.observe(el));
}

document.addEventListener("DOMContentLoaded", () => {
  initNavbar();
  initReveal();
  initCounters();
});
