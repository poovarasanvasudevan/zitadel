!function(){"use strict";const e=["client/client.b9ebf174.js","client/en.a7ea2c75.js","client/de.fd2e2b12.js","client/index.c96e95a5.js","client/[slug].f7f986cd.js","client/client.a26c42b2.js"].concat(["service-worker-index.html","base.css","fonts/ailerons/ailerons.otf","fonts/fira-mono/fira-mono-latin-400.woff2","fonts/roboto/roboto-latin-400.woff2","fonts/roboto/roboto-latin-400italic.woff2","fonts/roboto/roboto-latin-500.woff2","fonts/roboto/roboto-latin-500italic.woff2","icons/android-chrome-192x192.png","icons/android-chrome-512x512.png","icons/apple-touch-icon.png","icons/favicon-16x16.png","icons/favicon-32x32.png","icons/favicon.ico","icons/mstile-150x150.png","icons/safari-pinned-tab.svg","logos/zitadel-logo-oneline-darkdesign.svg","logos/zitadel-logo-solo-darkdesign.svg","manifest.json","prism.css"]),t=new Set(e);self.addEventListener("install",t=>{t.waitUntil(caches.open("cache1597753438716").then(t=>t.addAll(e)).then(()=>{self.skipWaiting()}))}),self.addEventListener("activate",e=>{e.waitUntil(caches.keys().then(async e=>{for(const t of e)"cache1597753438716"!==t&&await caches.delete(t);self.clients.claim()}))}),self.addEventListener("fetch",e=>{if("GET"!==e.request.method||e.request.headers.has("range"))return;const o=new URL(e.request.url);o.protocol.startsWith("http")&&(o.hostname===self.location.hostname&&o.port!==self.location.port||(o.host===self.location.host&&t.has(o.pathname)?e.respondWith(caches.match(e.request)):"only-if-cached"!==e.request.cache&&e.respondWith(caches.open("offline1597753438716").then(async t=>{try{const o=await fetch(e.request);return t.put(e.request,o.clone()),o}catch(o){const n=await t.match(e.request);if(n)return n;throw o}}))))})}();
