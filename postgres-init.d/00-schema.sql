create table if not exists sponsor_times
(
    "videoID"        varchar(12)                   not null,
    "startTime"      real                          not null,
    "endTime"        real                          not null,
    votes            bigint                        not null,
    locked           boolean     default false     not null,
    "incorrectVotes" integer                       not null,
    "UUID"           varchar(128) unique           not null,
    "userID"         varchar(128)                  not null,
    "timeSubmitted"  bigint                        not null,
    views            bigint                        not null,
    category         varchar(80) default 'sponsor' not null,
    "actionType"     varchar(30) default 'skip'    not null,
    service          varchar(80) default 'YouTube' not null,
    "videoDuration"  bigint      default 0         not null,
    hidden           boolean     default false     not null,
    reputation       integer     default 0         not null,
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
