FROM node:26.4.0-alpine AS builder

WORKDIR /build

COPY package*.json /build/
RUN npm install

COPY . /build/

# runner
FROM node:26.4.0-alpine AS runner

WORKDIR /app
COPY --from=builder /build /app/

EXPOSE 8083

CMD ["npm", "run", "app"]
