import os

# Read emails and alerts count from file
def read_gmail_addresses(file_path):
    email_alerts = []
    try:
        with open(file_path, 'r') as file:
            for line in file:
                email, count = line.strip().split(':')
                email_alerts.append((email, int(count)))
    except FileNotFoundError:
        pass
    return email_alerts

# Write a single email's alerts count to file
def update_gmail_alert_count(gmail_accounts, gmail_index):
    file_path = 'gmail_accounts.txt'
    with open(file_path, 'w') as file:
        for i, (gmail, count) in enumerate(gmail_accounts):
            if i == gmail_index:
                count += 1
            file.write(f"{gmail}:{count}\n")
    return count

# Function to initialize email alerts
def initialize_gmail_accounts():
    GMAIL_ACCOUNTS_FILE = os.getenv('GMAIL_ACCOUNTS_FILE', 'gmail_accounts.txt')
    email_alerts = read_gmail_addresses(GMAIL_ACCOUNTS_FILE)
    if not email_alerts:
        print("No emails found in the file.")
        return None, None

    current_email_index = 0
    return email_alerts, current_email_index