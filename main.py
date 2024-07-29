from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.service import Service
import time

def set_flight_alert(driver, origin, destination):
    time.sleep(2)  # Wait for the page to load
    try:
        # Click the cookies "Accept All" button
        accept_cookies_xpath = '//*[@id="yDmH0d"]/c-wiz/div/div/div/div[2]/div[1]/div[3]/div[1]/div[1]/form[2]/div/div/button/span'
        click_element(driver, accept_cookies_xpath)
        
        # Select one way flight
        trip_type_xpath = '//*[@id="yDmH0d"]/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[1]/div[1]/div[1]/div/div[1]/div[1]/div/div/div/div[1]/div'
        click_element(driver, trip_type_xpath)
        
        one_way_xpath = '//*[@id="yDmH0d"]/c-wiz[2]/div/div[2]/c-wiz/div[1]/c-wiz/div[2]/div[1]/div[1]/div[1]/div/div[1]/div[1]/div/div/div/div[2]/ul/li[2]'
        click_element(driver, one_way_xpath)
        
        # Enter the destination
        destination_input_xpath = '//*[@id="i23"]/div[4]/div/div/div[1]/div/div/input'
        enter_text(driver, destination_input_xpath, destination)
        
        search_button_xpath = '//*[@id="c2"]'  # XPath for the search button
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
        
def click_element(driver, xpath):
    """Clicks on the element specified by the given XPath."""
    wait_time = 1
    try:
        element = driver.find_element(By.XPATH, xpath)
        element.click()
        time.sleep(wait_time)  # Wait for the action to complete
        print(f"Element clicked successfully: {xpath}")
    except Exception as e:
        print(f"An error occurred while clicking element: {xpath}, Error: {e}")
        
def enter_text(driver, xpath, text):
    """Enters the given text into the element specified by the given XPath."""
    wait_time=1
    try:
        element = driver.find_element(By.XPATH, xpath)
        element.clear()  # Clear any pre-existing text
        element.send_keys(text)
        time.sleep(wait_time)  # Wait for the action to complete
        print(f"Text '{text}' entered successfully in element: {xpath}")
    except Exception as e:
        print(f"An error occurred while entering text: {xpath}, Error: {e}")

def main():
    # Setup ChromeDriver path
    service = Service('/opt/homebrew/Caskroom/chromedriver/127.0.6533.72/chromedriver-mac-arm64/chromedriver')
    driver = webdriver.Chrome(service=service)

    # Open Google Flights
    driver.get('https://www.google.com/travel/flights?gl=GB&hl=en-GB')

    routes = [('MAN', 'Paris')]

    for origin, destination in routes:
        set_flight_alert(driver, origin, destination)
        time.sleep(500)  # Wait a bit before setting up the next alert

    driver.quit()

if __name__ == "__main__":
    main()
