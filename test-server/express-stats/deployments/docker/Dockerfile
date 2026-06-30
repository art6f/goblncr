FROM node:26.4.0-alpine

WORKDIR /app

COPY package*.json .
RUN npm install
COPY . .

EXPOSE 8083

CMD ["npx", "tsx", "src/index.ts"]
