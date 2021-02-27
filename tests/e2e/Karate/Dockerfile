FROM maven:3-openjdk-8

COPY ./waitWebApp.sh /waitWebApp.sh
RUN sed -i 's/\r$//' /waitWebApp.sh
RUN chmod +x /waitWebApp.sh

WORKDIR /project
COPY . /project

RUN mvn clean install -Dmaven.test.skip=true

CMD /waitWebApp.sh && mvn test -Dtest=TestParallel