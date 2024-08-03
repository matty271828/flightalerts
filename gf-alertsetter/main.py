import time

from gmail import initialize_gmail_accounts, update_gmail_alert_count
from routes import routes
from page_interactions import accept_cookies, sign_in, populate_search_page, set_flight_alert, reset_search_page

from dotenv import load_dotenv
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
        
def main():
    # Load environment variables from .env file
    load_dotenv()
    
    # Read emails and initialize alerts count
    gmail_accounts, current_account_index = initialize_gmail_accounts()
    if gmail_accounts is None:
        return

    # Open Google Flights
    service = Service('/opt/homebrew/Caskroom/chromedriver/127.0.6533.72/chromedriver-mac-arm64/chromedriver')
    driver = webdriver.Chrome(service=service)
    driver.get('https://www.google.com/travel/flights?gl=GB&hl=en-GB')
    
    accept_cookies(driver)
    sign_in(driver, gmail_accounts[current_account_index][0])
    
    # Custom handling for the first route
    origin, destination = routes[0]
    populate_search_page(driver, origin, destination)     
    set_flight_alert(driver)

    for origin, destination in routes[1:len(routes)]:
        if gmail_accounts[current_account_index][1] >= 100:
            current_account_index += 1
            if current_account_index >= len(gmail_accounts):
                print("All emails have reached the alert limit.")
                break
            # Force a sign out
            driver.quit()
            time.sleep(3)
            driver = webdriver.Chrome(service=service)
            driver.get('https://www.google.com/travel/flights?gl=GB&hl=en-GB')
            # sign in using the next email address
            accept_cookies(driver)
            sign_in(driver, gmail_accounts[current_account_index][0])

        print(f"Attempting to set alert for: {origin} -> {destination} using email: {gmail_accounts[current_account_index][0]}")
        reset_search_page(driver)
        populate_search_page(driver, origin, destination)       
        set_flight_alert(driver)
        update_gmail_alert_count(gmail_accounts, current_account_index)
        
    time.sleep(2)
    driver.quit()

if __name__ == "__main__":
    main()
