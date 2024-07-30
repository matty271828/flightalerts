import time
import os

from element_interactions import click_element, enter_text

from dotenv import load_dotenv
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service

def populate_initial_search_page(driver, origin, destination):
    time.sleep(2)  # Wait for the page to load
    try:
        # Select one way flight
        trip_type_xpath = '//*[@id="yDmH0d"]/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[1]/div[1]/div[1]/div/div[1]/div[1]/div/div/div/div[1]/div'
        click_element(driver, trip_type_xpath)
        
        one_way_xpath = '//*[@id="yDmH0d"]/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[1]/div[1]/div[1]/div/div[1]/div[1]/div/div/div/div[2]/ul/li[2]'
        click_element(driver, one_way_xpath)
        
        # Enter the destination
        destination_input_xpath = '//*[@id="i23"]/div[4]/div/div/div[1]/div/div/input'
        enter_text(driver, destination_input_xpath, destination)
        
        search_button_xpath = '/html/body/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[1]/div[1]/div[1]/div/div[2]/div[1]/div[6]/div[3]/ul/li[1]'  # XPath for the search button
        click_element(driver, search_button_xpath)

        # Enter the origin
        origin_input_xpath = '//*[@id="i23"]/div[1]/div/div/div[1]/div/div/input'
        enter_text(driver, origin_input_xpath, origin)
        
        origin_search_result_xpath = '/html/body/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[1]/div[1]/div[1]/div/div[2]/div[1]/div[6]/div[3]/ul/li[1]'  # XPath for the origin search result
        click_element(driver, origin_search_result_xpath)
        
        # Select a date - doesnt matter what it is as we will be later selecting any date for alerts
        calendar_input_xpath = '//*[@id="yDmH0d"]/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[1]/div[1]/div[1]/div/div[2]/div[2]/div/div/div[1]/div/div/div[1]/div/div[1]/div/input'
        click_element(driver, calendar_input_xpath)
        
        date_input_xpath = '//*[@id="ow79"]/div[2]/div/div[2]/div[2]/div/div/div[1]/div/div[2]/div[3]/div[1]/div[4]/div/div[2]'
        click_element(driver, date_input_xpath)
        
        done_button_xpath = '//*[@id="ow79"]/div[2]/div/div[3]/div[3]/div/button'
        click_element(driver, done_button_xpath)
        
        # click on search 
        search_button_xpath = '//*[@id="yDmH0d"]/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[1]/div[1]/div[2]/div/button'
        click_element(driver, search_button_xpath)
    
    except Exception as e:
        print(f"An error occurred: {e}")
        
def sign_in(driver):
    try:
        sign_in_button_xpath = '//*[@id="gb"]/div[2]/div[3]/div[1]/a'
        click_element(driver, sign_in_button_xpath)
        
        email_input_xpath = '//*[@id="identifierId"]'
        email = os.getenv('EMAIL')
        enter_text(driver, email_input_xpath, email)
        
        next_button_xpath = '//*[@id="identifierNext"]/div/button'
        click_element(driver, next_button_xpath)
        
        password_input_xpath = '//*[@id="password"]/div[1]/div/div[1]/input'
        password = os.getenv('PASSWORD')
        enter_text(driver, password_input_xpath, password)
        
        next_button_xpath = '//*[@id="passwordNext"]/div/button'
        click_element(driver, next_button_xpath)
    
    except Exception as e:
        print(f"An error occurred: {e}")
        
def set_flight_alert(driver):
    try:
        set_alert_xpath = '/html/body/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[2]/div[2]/div[1]/div/div[1]/label[2]/span[2]/span[2]/button'
        
        # Find the button element using XPath
        button_element = driver.find_element(By.XPATH, set_alert_xpath)
        
        # Check if the aria-checked attribute is set to 'false'
        aria_checked = button_element.get_attribute("aria-checked")
        
        if aria_checked == 'false':
            # Click the button if aria-checked is 'false'
            click_element(driver, set_alert_xpath)
        else:
            print("Alert is already set.")       
        
    except Exception as e:
        print(f"An error occurred: {e}")
        
def accept_cookies(driver):
    try:
        # Click the cookies "Accept All" button
        accept_cookies_xpath = '//*[@id="yDmH0d"]/c-wiz/div/div/div/div[2]/div[1]/div[3]/div[1]/div[1]/form[2]/div/div/button/span'
        click_element(driver, accept_cookies_xpath)
    except Exception as e:
        print(f"An error occurred: {e}")
        
        
def main():
    # Load environment variables from .env file
    load_dotenv()

    # Setup ChromeDriver path
    service = Service('/opt/homebrew/Caskroom/chromedriver/127.0.6533.72/chromedriver-mac-arm64/chromedriver')
    driver = webdriver.Chrome(service=service)

    # Open Google Flights
    driver.get('https://www.google.com/travel/flights?gl=GB&hl=en-GB')

    routes = [('MAN', 'Paris')]

    for origin, destination in routes:
        accept_cookies(driver)
        
        sign_in(driver)
        
        populate_initial_search_page(driver, origin, destination)     
        
        set_flight_alert(driver)
        
        time.sleep(500)  # Wait a bit before setting up the next alert

    driver.quit()

if __name__ == "__main__":
    main()
