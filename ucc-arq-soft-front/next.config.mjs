/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: [
      "upload.wikimedia.org",
      "encrypted-tbn0.gstatic.com",
      "i.imgflip.com",
      "res.cloudinary.com",
      "ih1.redbubble.net",
    ],
  },
  // Disable static export / prerendering of app routes to avoid
  // "Unsupported Server Component type" errors with client-only pages.
  output: "standalone",
  trailingSlash: false,
};

export default nextConfig;
