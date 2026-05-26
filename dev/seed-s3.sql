CREATE DATABASE IF NOT EXISTS analytics_s3;

CREATE TABLE IF NOT EXISTS analytics_s3.events
(
    event_id UUID DEFAULT generateUUIDv4(),
    event_time DateTime DEFAULT now(),
    event_type LowCardinality(String),
    user_id UInt64,
    session_id String,
    page_url String,
    country LowCardinality(String),
    city LowCardinality(String),
    browser LowCardinality(String),
    os LowCardinality(String),
    device LowCardinality(String),
    duration_ms UInt32,
    bytes_sent UInt64,
    status_code UInt16
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(event_time)
ORDER BY (event_time, user_id, event_type)
SETTINGS disk = 'disk_s3', index_granularity = 8192;

INSERT INTO analytics_s3.events (event_time, event_type, user_id, session_id, page_url, country, city, browser, os, device, duration_ms, bytes_sent, status_code)
SELECT
    now() - intDiv(number, 10) * 60,
    ['page_view', 'click', 'api_call', 'error', 'purchase'][rand32() % 5 + 1],
    rand64() % 100000,
    concat('session_', toString(rand64() % 10000)),
    concat('/page/', toString(rand64() % 500)),
    ['US', 'UK', 'DE', 'FR', 'JP', 'AU', 'CA', 'BR', 'IN', 'NL'][rand32() % 10 + 1],
    ['New York', 'London', 'Berlin', 'Paris', 'Tokyo', 'Sydney', 'Toronto', 'Sao Paulo', 'Mumbai', 'Amsterdam'][rand32() % 10 + 1],
    ['Chrome', 'Firefox', 'Safari', 'Edge'][rand32() % 4 + 1],
    ['Windows', 'macOS', 'Linux', 'iOS', 'Android'][rand32() % 5 + 1],
    ['desktop', 'mobile', 'tablet'][rand32() % 3 + 1],
    rand32() % 5000,
    rand64() % 1000000,
    [200, 200, 200, 200, 404, 500, 301][rand32() % 7 + 1]
FROM numbers(50000);

SYSTEM FLUSH LOGS;

SELECT count(*) FROM analytics_s3.events;

SELECT
    country,
    city,
    count() AS cnt,
    avg(duration_ms) AS avg_duration,
    sum(bytes_sent) AS total_bytes
FROM analytics_s3.events
GROUP BY country, city
ORDER BY cnt DESC
LIMIT 20;

SELECT
    browser,
    os,
    device,
    count() AS cnt,
    quantile(0.95)(duration_ms) AS p95_duration
FROM analytics_s3.events
WHERE event_type IN ('page_view', 'click')
GROUP BY browser, os, device
ORDER BY cnt DESC
LIMIT 15;

SELECT
    toStartOfHour(event_time) AS hour,
    event_type,
    count() AS cnt,
    uniqExact(user_id) AS unique_users,
    sum(bytes_sent) AS total_bytes
FROM analytics_s3.events
GROUP BY hour, event_type
ORDER BY hour DESC, cnt DESC
LIMIT 50;

SYSTEM FLUSH LOGS;
