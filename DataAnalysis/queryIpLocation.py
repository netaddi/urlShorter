# imports
import pymysql
import urllib.request
import json

dbConfigStr = open("db.json", 'r').read()
dbConfig    = json.loads(dbConfigStr)
dbConfig

# get db cursor
db = pymysql.connect(host     = dbConfig["host"],
                     user     = dbConfig["user"],
                     passwd   = dbConfig["password"],
                     db       = dbConfig["db"],
                     charset  = 'utf8')



cursor = db.cursor()


# add ip with unknown addresses

ipListUpdateQuery = """INSERT INTO ip_info (ip)
  SELECT DISTINCT forwarded_ip
  FROM access_log
  WHERE forwarded_ip
        NOT IN (
          SELECT ip
          FROM ip_info
        )"""
try:
    cursor.execute(ipListUpdateQuery)
    db.commit()
except:
    db.rollback()


# query ip addresses

query  = "SELECT ip FROM ip_info WHERE country IS NULL"
cursor.execute(query)
results = cursor.fetchall()


for row in results:
    print(row[0])
    ip = row[0]
    locationQueryUrl = "http://ip.taobao.com/service/getIpInfo.php?ip=" + ip
    locationData = urllib.request.urlopen(locationQueryUrl).read().decode('utf-8')
    locationData = json.loads(locationData)["data"]
    print(locationData)

    updateQuery = """UPDATE ip_info
                       SET country = '%s',
                           province = '%s',
                           city = '%s',
                           isp = '%s'
                       WHERE ip = '%s'""" % (locationData["country"],
                                           locationData["region"],
                                           locationData["city"],
                                           locationData["isp"],
                                           ip)
    print(updateQuery)

    cursor.execute(updateQuery.encode('utf8'))

try:
    db.commit()
except:
    db.rollback()
db.close()
