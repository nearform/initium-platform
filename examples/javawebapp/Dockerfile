FROM adoptopenjdk:11-jdk-hotspot
WORKDIR /app
COPY target/web-application-1.0.0.jar /app/web-application.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "/app/web-application.jar"]
