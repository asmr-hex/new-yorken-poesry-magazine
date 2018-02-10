FROM node:9.4.0-alpine

WORKDIR /client

COPY package*.json ./

RUN npm install -qy

COPY . .

EXPOSE 3000

CMD ["npm", "start"]
