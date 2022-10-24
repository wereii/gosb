create table if not exists sponsor_times
(
    "videoID"        text                          not null,
    "startTime"      real                          not null,
    "endTime"        real                          not null,
    votes            bigint                        not null,
    locked           smallint    default 0         not null,
    "incorrectVotes" integer                       not null,
    "UUID"           varchar(128) unique           not null,
    "userID"         text                          not null,
    "timeSubmitted"  bigint                        not null,
    views            bigint                        not null,
    category         text        default 'sponsor' not null,
    "actionType"     text        default 'skip'    not null,
    service          text        default 'YouTube' not null,
    "videoDuration"  real        default 0         not null,
    hidden           smallint    default 0         not null,
    reputation       real        default 0         not null,
    "shadowHidden"   boolean     default false     not null,
    "hashedVideoID"  varchar(64) default ''        not null,
    "userAgent"      text        default ''        not null,
    description      text        default ''        not null
);

-- create index if not exists "idx_sponsorTime_timeSubmitted"
--    on sponsor_times ("timeSubmitted");

--create index if not exists "idx_sponsorTime_startTime"
--    on sponsor_times ("startTime");

create index if not exists "idx_sponsorTimes_skipSegments"
    on sponsor_times ("hashedVideoID" varchar_pattern_ops, "votes", "category", "startTime");

-- Silently drop invalid stuff, it's actually easier to do this here then on the \copy side
CREATE OR REPLACE FUNCTION cleanup_sponsor_time_insert()
    RETURNS trigger AS
$func$
BEGIN
    -- If not alphanumeric_-
    IF NEW."videoID" !~* '^[a-z0-9_\-]*' THEN
        RAISE NOTICE 'Non-alphanumeric VideoID "%" ', NEW."VideoID";
    ELSEIF NEW."UUID" !~* '^[a-f0-9]*' THEN
        -- If not hexadecimal
        RAISE NOTICE 'Non-hexadecimal UUID "%" ', NEW."UUID";
    ELSEIF NEW."userID" !~* '^[a-f0-9]*' THEN
        -- If not hexadecimal
        RAISE NOTICE 'Non-hexadecimal userID "%" ', NEW."userID";
    ELSEIF NEW."hashedVideoID" !~* '^[a-f0-9]*' THEN
        -- If not hexadecimal
        RAISE NOTICE 'Non-hexadecimal hashedVideoID "%" ', NEW."hashedVideoID";
    ELSE
        RETURN NEW;
    END IF;

    -- drop silently as to not break csv copy/import
    RETURN NULL;
END;
$func$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER cleanup_trigger
    BEFORE INSERT
    ON sponsor_times
    FOR EACH ROW
EXECUTE FUNCTION cleanup_sponsor_time_insert();
