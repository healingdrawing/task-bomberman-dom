FROM node:19.6.0 as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./ .
RUN npm run build

FROM nginx:1.21.0 as production-stage
RUN mkdir /app
COPY --from=build-stage /app/dist /app/dist
COPY --from=build-stage /app/assets /app/assets
COPY --from=build-stage /app/index.html /app/index.html
COPY nginx.conf /etc/nginx/nginx.conf
