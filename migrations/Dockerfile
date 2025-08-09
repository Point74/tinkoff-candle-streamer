FROM liquibase/liquibase:4.25.0
USER 2000
ADD /db /db
CMD ["sh", "-c", "docker-entrypoint.sh --url=${DB_URL} --username=${DB_USERNAME} --password=${DB_PASSWORD} --classpath=/ --changeLogFile=db/db-changelog.xml update"]