#include <iostream>
#include <fstream>
#include <vector>
#include <map>
#include <sstream>
#include <cmath>
#include <algorithm>
#include <ctime>
#include <iomanip>

struct Object {
    std::string name;
    double x;
    double y;
    std::string typ;
    time_t time; // Используем time_t для хранения времени
};

double calculateDistance(double x, double y) {
    return std::sqrt(x * x + y * y);
}

void sortAndWriteByDistance(const std::vector<Object>& objects, const std::string& fil) {
    std::vector<Object> upTo100, upTo1000, upTo10000, upTo10000plus;

    for (const Object& obj : objects) {
        double dist = calculateDistance(obj.x, obj.y);
        if (dist < 100) {
            upTo100.push_back(obj);
        } else if (dist < 1000) {
            upTo1000.push_back(obj);
        } else if (dist < 10000) {
            upTo10000.push_back(obj);
        } else {
            upTo10000plus.push_back(obj);
        }
    }

    auto distanceComparator = [](const Object& obj1, const Object& obj2) {
        double dist1 = calculateDistance(obj1.x, obj1.y);
        double dist2 = calculateDistance(obj2.x, obj2.y);
        return dist1 < dist2;
    };

    std::sort(upTo100.begin(), upTo100.end(), distanceComparator);
    std::sort(upTo1000.begin(), upTo1000.end(), distanceComparator);
    std::sort(upTo10000.begin(), upTo10000.end(), distanceComparator);

    std::ofstream file(fil);
    if (!file.is_open()) {
        std::cerr << "Ошибка создания файла" << std::endl;
        return;
    }

    auto writeGroup = [&file](const std::vector<Object>& group, const std::string& label) {
        if (!group.empty()) {
            file << "=== Расстояние " << label << " ===" << std::endl;
            for (const Object& obj : group) {
                file << obj.name << ", Расстояние: " << calculateDistance(obj.x, obj.y) << ", Тип: " << obj.typ << ", Время: " << obj.time << std::endl;
            }
        }
    };

    writeGroup(upTo100, "до 100");
    writeGroup(upTo1000, "до 1000");
    writeGroup(upTo10000, "до 10000");

    std::cout << "Результаты записаны в файл." << std::endl;
}
void grname(const std::vector<Object>& objects, const std::string& fil) {
    std::map<std::string, std::vector<Object>> groups;

    for (const Object& obj : objects) {
        std::string name = obj.name;
        std::string groupName;

        bool isCyrillic = false;
        for (char c : name) {
            if (std::iswalpha(c) && std::iswalpha(static_cast<wint_t>(c))) {
                isCyrillic = true;
                break;
            }
        }

        if (isCyrillic) {
        } else {
            groupName = "";
        }

        groups[groupName].push_back(obj);
    }

    std::ofstream file(fil);
    if (!file.is_open()) {
        std::cerr << "Ошибка создания файла" << std::endl;
        return;
    }

    for (const auto& group : groups) {
        if (!group.first.empty()) {
            file << "Группа " << group.first << ":" << std::endl;
        }
        for (const Object& obj : group.second) {
            file << obj.name << ", Расстояние: " << calculateDistance(obj.x, obj.y) << ", Тип: " << obj.typ << ", Время: " << obj.time << std::endl;
        }
    }

    std::cout << "Результаты группировки записаны в файл." << std::endl;
}


void grtime(const std::vector<Object>& objects, const std::string& fil) {
    std::map<std::string, std::vector<Object>> groups;

    for (const Object& obj : objects) {
        time_t createTime = obj.time;

        std::string groupName;
        time_t currentTime = std::time(nullptr);

        if (createTime > currentTime - 24 * 3600) {
            groupName = "Сегодня";
        } else if (createTime > currentTime - 2 * 24 * 3600) {
            groupName = "Вчера";
        } else if (createTime > currentTime - 7 * 24 * 3600) {
            groupName = "На этой неделе";
        } else if (createTime > currentTime - 30 * 24 * 3600) {
            groupName = "В этом месяце";
        } else if (createTime > currentTime - 365 * 24 * 3600) {
            groupName = "В этом году";
        } else {
            groupName = "Ранее";
        }

        groups[groupName].push_back(obj);
    }

    std::ofstream file(fil);
    if (!file.is_open()) {
        std::cerr << "Ошибка создания файла" << std::endl;
        return;
    }

    for (const auto& group : groups) {
        if (!group.first.empty()) {
            file << "Группа " << group.first << ":" << std::endl;
        }
        for (const Object& obj : group.second) {
            char timeBuffer[20];
            strftime(timeBuffer, sizeof(timeBuffer), "%Y-%m-%d %H:%M:%S", localtime(&obj.time));
            file << obj.name << ", Расстояние: " << calculateDistance(obj.x, obj.y) << ", Тип: " << obj.typ << ", Время: " << timeBuffer << std::endl;
        }
    }

    std::cout << "Результаты группировки по времени записаны в файл." << std::endl;
}

void grtype(const std::vector<Object>& objects, const std::string& fil) {
    int minCount;
    std::cout << "Введите минимальное количество объектов в группе: ";
    std::cin >> minCount;

    std::map<std::string, std::vector<Object>> groups;

    for (const Object& obj : objects) {
        groups[obj.typ].push_back(obj);
    }

    for (auto it = groups.begin(); it != groups.end();) {
        if (it->second.size() < minCount) {
            it = groups.erase(it);
        } else {
            ++it;
        }
    }

    for (auto& group : groups) {
        std::sort(group.second.begin(), group.second.end(), [](const Object& obj1, const Object& obj2) {
            return obj1.name < obj2.name;
        });
    }

    std::ofstream file(fil);
    if (!file.is_open()) {
        std::cerr << "Ошибка создания файла" << std::endl;
        return;
    }

    for (const auto& group : groups) {
        file << "Группа " << group.first << ":" << std::endl;
        for (const Object& obj : group.second) {
            file << obj.name << ", Расстояние: " << calculateDistance(obj.x, obj.y) << ", Тип: " << obj.typ << ", Время: " << obj.time << std::endl;
        }
    }

    std::cout << "Результаты группировки по типу записаны в файл." << std::endl;
}

int main() {
    std::cout << "Введите путь к файлу: ";
    std::string fil;
    std::cin >> fil;

    std::ifstream file(fil);
    if (!file.is_open()) {
        std::cerr << "Ошибка открытия файла" << std::endl;
        return 1;
    }

    std::vector<Object> objects;
    std::string line;

    while (std::getline(file, line)) {
        std::istringstream iss(line);
        Object obj;
        if (iss >> obj.name >> obj.x >> obj.y >> obj.typ >> obj.time) {
            objects.push_back(obj);
        }
    }

    file.close();

    if (objects.empty()) {
        std::cerr << "Нет данных для обработки" << std::endl;
        return 1;
    }

    std::string choice;
    std::cout << "Выберите действие:" << std::endl;
    std::cout << "1. Группировка по расстоянию." << std::endl;
    std::cout << "2. Группировка по имени." << std::endl;
    std::cout << "3. Группировка по времени создания." << std::endl;
    std::cout << "4. Группировка по типу." << std::endl;
    std::cin >> choice;

    if (choice == "1") {
        sortAndWriteByDistance(objects, fil);
    } else if (choice == "2") {
        grname(objects, fil);
    } else if (choice == "3") {
        grtime(objects, fil);
    } else if (choice == "4") {
        grtype(objects, fil);
    } else {
        std::cerr << "Неверный выбор" << std::endl;
    }

    return 0;
}
