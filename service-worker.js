!function(){"use strict";const t=["client/client.7df290cc.js","client/en.e16f1f3f.js","client/de.16573860.js","client/index.62e25283.js","client/[slug].7a37edb5.js","client/client.853dfc12.js"].concat(["service-worker-index.html","base.css","fonts/ailerons/ailerons.otf","fonts/fira-mono/fira-mono-latin-400.woff2","fonts/roboto/roboto-latin-400.woff2","fonts/roboto/roboto-latin-400italic.woff2","fonts/roboto/roboto-latin-500.woff2","fonts/roboto/roboto-latin-500italic.woff2","icons/android-chrome-192x192.png","icons/android-chrome-512x512.png","icons/apple-touch-icon.png","icons/favicon-16x16.png","icons/favicon-32x32.png","icons/favicon.ico","icons/mstile-150x150.png","icons/safari-pinned-tab.svg","logos/zitadel-logo-dark@2x.png","logos/zitadel-logo-light.svg","logos/zitadel-logo-solo-darkdesign.svg","manifest.json","prism.css"]),o=new Set(t);self.addEventListener("install",o=>{o.waitUntil(caches.open("cache1600761921983").then(o=>o.addAll(t)).then(()=>{self.skipWaiting()}))}),self.addEventListener("activate",t=>{t.waitUntil(caches.keys().then(async t=>{for(const o of t)"cache1600761921983"!==o&&await caches.delete(o);self.clients.claim()}))}),self.addEventListener("fetch",t=>{if("GET"!==t.request.method||t.request.headers.has("range"))return;const e=new URL(t.request.url);e.protocol.startsWith("http")&&(e.hostname===self.location.hostname&&e.port!==self.location.port||(e.host===self.location.host&&o.has(e.pathname)?t.respondWith(caches.match(t.request)):"only-if-cached"!==t.request.cache&&t.respondWith(caches.open("offline1600761921983").then(async o=>{try{const e=await fetch(t.request);return o.put(t.request,e.clone()),e}catch(e){const n=await o.match(t.request);if(n)return n;throw e}}))))})}();
