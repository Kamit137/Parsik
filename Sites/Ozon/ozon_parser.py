import json
import time
import undetected_chromedriver as uc
from bs4 import BeautifulSoup
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
import sys


url = "https:///www.ozon.ru/product/sredstvo-dlya-mytya-posudy-synergetic-detskih-igrushek-c-aromatom-granata-5-l-antibakterialnoe-1436053626/"

# url = sys.argv[1]

def collect_product_info(driver, url):
    product_id = WebDriverWait(driver, 10).until(EC.presence_of_element_located((By.XPATH, '//div[contains(text(), "Артикул: ")]'))).text.split('Артикул: ')[1]
    image_url = None
    page_source = str(driver.page_source)
    soup = BeautifulSoup(page_source, 'lxml')

    with open(f'product_{product_id}.html', 'w') as file:
        file.write(page_source)

    product_name = soup.find('div', attrs={"data-widget": 'webProductHeading'}).find(
        'h1').text.strip().replace('\t', '').replace('\n', ' ')
    try:
            ozon_card_price_element = soup.find(
                'span', string="c Ozon Картой").parent.find('div').find('span')
            product_ozon_card_price = ozon_card_price_element.text.strip(
            ) if ozon_card_price_element else ''

            price_element = soup.find(
                'span', string="без Ozon Карты").parent.parent.find('div').find_all('span')

            product_discount_price = price_element[0].text.strip(
            ) if price_element[0] else ''
            product_base_price = price_element[1].text.strip(
            ) if price_element[1] is not None else ''
    except:
        product_ozon_card_price = None
        product_discount_price = None


        # product price
        try:
            ozon_card_price_element = soup.find(
                'span', string="c Ozon Картой").parent.find('div').find('span')
        except AttributeError:
            card_price_div = soup.find(
                'div', attrs={"data-widget": "webPrice"}).findAll('span')


            product_discount_price = card_price_div[1].text.strip()


    try:
        img_element = WebDriverWait(driver, 10).until(EC.presence_of_element_located((By.CSS_SELECTOR, 'div[data-widget="webGallery"] img, img.zoom-image')))

        # Получаем URL маленького изображения
        small_image_url = img_element.get_attribute('src')
        # Преобразуем в URL большого изображения
        if small_image_url:
            # Заменяем параметры размера на максимальный
            image_url = (small_image_url.replace('/wc50/', '/wc1000/').replace('/wc100/', '/wc1000/').replace('/wc250/', '/wc1000/').split('?')[0])  # Удаляем параметры URL после ?
    except:
        image_url = None

    product = {
        'Name': product_name,
        # 'product_id': product_id,
        'Price': product_ozon_card_price.replace("\u2009","").replace("\u2009","").replace("₽",""),

        'imageURL': image_url,

        # 'product_discount_price': product_discount_price.replace("\u2009","").replace("\u2009","").replace("₽",""),

    }


    return product

driver = uc.Chrome()

driver.get(url)


products_data = []


data = collect_product_info(driver=driver, url=url)
print(json.dumps(data, ensure_ascii=False, indent=2))

driver.close()
driver.quit()
