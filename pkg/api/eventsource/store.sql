CREATE table events_store(
    resourceID VARCHAR(36),
    aggregate bool,
    version int,
    state VARCHAR(50) NOT NULL,
    Content blob,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (resourceID, version)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
