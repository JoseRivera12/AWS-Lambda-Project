import csv
import sys
import random
from datetime import datetime, timedelta

def random_date_2024():
    start_date = datetime(2024, 1, 1)
    end_date = datetime(2024, 12, 31)
    return start_date + timedelta(seconds=random.randint(0, int((end_date - start_date).total_seconds())))


def random_transaction_amount():
    return round(random.uniform(-1000, 1000), 2)

def random_description():
    descriptions = ["School expences", "Restaurant bill", "Grocery shopping", "Videogames", "Utilities bill", "Rent payment", "Salary", "Online purchase", "Medical expenses", "Refunds"]
    return random.choice(descriptions)

def main():
    try:
        file_name = sys.argv[1]
        num_transactions = int(sys.argv[2])
    except:
        print("Run as follow: python csv_generator user_transactions 100000")
        return 
    transactions = []
    for i in range(num_transactions):
        transaction_amount = random_transaction_amount()
        transaction = {
            "Id": i,
            "Date": random_date_2024().strftime("%Y-%m-%dT%H:%M:%S"),
            "Transaction": transaction_amount if transaction_amount < 0 else f"+{transaction_amount}",
            "Currency": "USD",
            "Description": random_description()
        }
        transactions.append(transaction)

    with open(f"{file_name}.csv", "w", newline="") as csvfile:
        fieldnames = ["Id", "Date", "Transaction", "Currency", "Description"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        for transaction in transactions:
            writer.writerow(transaction)

    print(f"{file_name}.csv with {num_transactions} transactions has been generated.")

if __name__ == "__main__":
    main()