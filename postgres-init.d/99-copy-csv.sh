if [ -e "/db-dumps" ] && [ -f "/db-dumps/sponsorTimes.csv" ]; then
  psql -d 'sb' -c "\COPY sponsor_times FROM '/db-dumps/sponsorTimes.csv' DELIMITER ',' CSV HEADER;"
else
  echo "Skipping import/copy, db-dumps folder or sponsorTimes.csv file not found"
fi
