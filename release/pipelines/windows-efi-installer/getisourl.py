#!/usr/bin/python

import sys

from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.select import Select

options = Options()
options.add_argument("--headless")

driver = webdriver.Chrome(options=options)
driver.implicitly_wait(10)

try:
    driver.get("https://www.microsoft.com/en-us/software-download/windows11")

    select_object = Select(driver.find_element(By.ID, "product-edition"))
    select_object.select_by_index(1)

    button_element = driver.find_element(By.ID, "submit-product-edition")
    button_element.click()

    select_object = Select(driver.find_element(By.ID, "product-languages"))
    select_object.select_by_visible_text("English (United States)")

    button_element = driver.find_element(By.ID, "submit-sku")
    button_element.click()

    download_links = driver.find_elements(By.PARTIAL_LINK_TEXT, "64-bit")
    if len(download_links) != 1:
        raise ValueError(
            f"expected exactly one 64-bit download link, found {len(download_links)}"
        )
    href = download_links[0].get_attribute("href")
    if not href:
        raise ValueError("download link has no href attribute")
    print(href)
except Exception as e:
    print(f"Failed to obtain the download URL: {e}", file=sys.stderr)
    sys.exit(1)
finally:
    driver.quit()
