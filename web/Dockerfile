FROM node:22-alpine AS base
WORKDIR /home/node/app
COPY package*.json ./
RUN npm install
COPY . .

FROM base AS development
EXPOSE 5173
CMD ["npm", "run", "dev"]

FROM base AS builder
RUN npm run build

FROM node:22-alpine AS production
WORKDIR /home/node/app
COPY --from=builder /home/node/app/dist ./dist
RUN npm install -g serve
EXPOSE 5173
CMD ["serve", "-s", "dist", "-l", "5173"]
