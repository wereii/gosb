create table if not exists "sponsorTimes"
(
    "videoID"        varchar(12)                   not null,
    "startTime"      real                          not null,
    "endTime"        real                          not null,
    votes            bigint                        not null,
    locked           bool        default false     not null,
    "incorrectVotes" integer                       not null,
    "UUID"           varchar(36)                   not null,
    "userID"         varchar(36)                   not null,
    "timeSubmitted"  bigint                        not null,
    views            bigint                        not null,
    category         varchar(80) default 'sponsor' not null,
    "actionType"     text        default 'skip'    not null,
    service          varchar(80) default 'YouTube' not null,
    "videoDuration"  bigint      default 0         not null,
    hidden           bool        default false     not null,
    reputation       integer     default 0         not null,
    "shadowHidden"   bool        default false     not null,
    "hashedVideoID"  varchar(64) default ''::text  not null,
    "userAgent"      text        default ''::text  not null,
    description      text        default ''::text  not null
);

create index if not exists "sponsorTime_timeSubmitted"
    on "sponsorTimes" ("timeSubmitted");

create index if not exists "sponsorTime_userID"
    on "sponsorTimes" ("userID");

create index if not exists "sponsorTimes_UUID"
    on "sponsorTimes" ("UUID");

create index if not exists "sponsorTimes_hashedVideoID"
    on "sponsorTimes" (service, "hashedVideoID" text_pattern_ops, "startTime");

create index if not exists "sponsorTimes_videoID"
    on "sponsorTimes" (service, "videoID", "startTime");

create index if not exists "sponsorTimes_videoID_category"
    on "sponsorTimes" ("videoID", category);
