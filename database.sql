DROP TABLE IF EXISTS management_log;
DROP TABLE IF EXISTS management_status;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS servers;
DROP TABLE IF EXISTS clusters;
DROP TABLE IF EXISTS sites;

CREATE TABLE sites (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE clusters (
    id INT PRIMARY KEY AUTO_INCREMENT,
    site INT NOT NULL,
    name VARCHAR(50) UNIQUE NOT NULL,
    FOREIGN KEY (site) REFERENCES sites(id)
);

CREATE TABLE servers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    site INT NOT NULL,
    cluster INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    serial VARCHAR(255) NOT NULL,
    ip VARCHAR(255) NOT NULL,
    FOREIGN KEY (site) REFERENCES sites(id),
    FOREIGN KEY (cluster) REFERENCES clusters(id),
    UNIQUE(site, cluster, name)
);

CREATE TABLE services (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE management_status (
    id INT PRIMARY KEY AUTO_INCREMENT,
    server INT NOT NULL,
    service INT NOT NULL,
    managed BOOLEAN NOT NULL,
    FOREIGN KEY (server) REFERENCES servers(id),
    FOREIGN KEY (service) REFERENCES services(id),
    UNIQUE(server, service)
);

CREATE TABLE management_log (
    id INT PRIMARY KEY AUTO_INCREMENT,
    dt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    managed INT NOT NULL,
    status BOOLEAN NOT NULL,
    FOREIGN KEY (managed) REFERENCES management_status(id)
);

DELIMITER //
CREATE TRIGGER management_status_log AFTER UPDATE ON management_status
FOR EACH ROW
BEGIN
	INSERT INTO management_log(managed, status) VALUES (NEW.id, NEW.managed);
END;
//
DELIMITER ;

INSERT INTO services (name) VALUES ('alerting');
INSERT INTO services (name) VALUES ('change-control');
