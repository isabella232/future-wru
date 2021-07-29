import csv
with open('testuser.csv', 'w', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(["id", "name", "mail", "org", "scopes", "twitter", "github", "oidc"])
    writer.writerow(["user1", "test user1", "testuser1@example.com", "R&D", "admin,user,org:rd", "testuser1", "testuser1", "testuser01@example.com"])
    writer.writerow(["user2", "test user2", "testuser2@example.com", "R&D", "admin,user,org:rd", "testuser2", "testuser2", "testuser02@example.com"])
    writer.writerow(["user2", "test user3", "testuser3@example.com", "R&D", "admin,user,org:rd", "shibu_jp", "shibukawa", "testuser03@example.com"])
