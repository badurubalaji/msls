-- Migration: 000002_uuid_v7.up.sql
-- Description: Create UUID v7 generation function
-- UUID v7 provides time-ordered UUIDs for better database performance

CREATE OR REPLACE FUNCTION uuid_generate_v7()
RETURNS uuid AS $$
DECLARE
    -- Current timestamp in milliseconds
    unix_ts_ms BIGINT;
    -- UUID bytes
    uuid_bytes BYTEA;
BEGIN
    -- Get current timestamp in milliseconds since Unix epoch
    unix_ts_ms := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT;

    -- Build 16 bytes for UUID
    -- First 6 bytes: timestamp (48 bits)
    -- Bytes 7-8: version (4 bits) + random (12 bits)
    -- Bytes 9-16: variant (2 bits) + random (62 bits)
    uuid_bytes := SET_BYTE(
        SET_BYTE(
            SET_BYTE(
                SET_BYTE(
                    SET_BYTE(
                        SET_BYTE(
                            gen_random_bytes(16),
                            0, ((unix_ts_ms >> 40) & 255)::INT
                        ),
                        1, ((unix_ts_ms >> 32) & 255)::INT
                    ),
                    2, ((unix_ts_ms >> 24) & 255)::INT
                ),
                3, ((unix_ts_ms >> 16) & 255)::INT
            ),
            4, ((unix_ts_ms >> 8) & 255)::INT
        ),
        5, (unix_ts_ms & 255)::INT
    );

    -- Set version to 7 (0111 in binary) - byte 6, high nibble
    uuid_bytes := SET_BYTE(uuid_bytes, 6, (GET_BYTE(uuid_bytes, 6) & 15) | 112);

    -- Set variant to RFC 4122 (10xx in binary) - byte 8, high 2 bits
    uuid_bytes := SET_BYTE(uuid_bytes, 8, (GET_BYTE(uuid_bytes, 8) & 63) | 128);

    RETURN ENCODE(uuid_bytes, 'hex')::uuid;
END;
$$ LANGUAGE plpgsql VOLATILE;

-- Add a comment for documentation
COMMENT ON FUNCTION uuid_generate_v7() IS 'Generates a UUID v7 (time-ordered) as per RFC 9562';
