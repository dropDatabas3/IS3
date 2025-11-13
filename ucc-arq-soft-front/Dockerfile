FROM node:18-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --silent --prefer-offline

FROM node:18-alpine AS builder
WORKDIR /app
# Allow passing NEXT_PUBLIC_API_URL at build time so Next.js embeds it correctly
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules
EXPOSE 3000
# Runtime config entrypoint: generate public/runtime-config.js from env
COPY docker-entrypoint.sh ./docker-entrypoint.sh
# Normalize line endings and ensure executable bit (Windows -> Linux)
RUN sed -i 's/\r$//' /app/docker-entrypoint.sh && chmod +x /app/docker-entrypoint.sh
ENTRYPOINT ["sh", "/app/docker-entrypoint.sh"]
CMD ["npm", "start"]