if [ -e "/db-dumps" ] && [ -f "/db-dumps/sponsorTimes.csv" ]; then
  psql -c "\COPY sponsorTimes FROM '/db-dumps/sponsorTimes.csv' DELIMITER ',' CSV HEADER;" || true
else
  echo "Skipping import/copy, db-dumps folder or sponsorTimes.csv file not found"
fi
