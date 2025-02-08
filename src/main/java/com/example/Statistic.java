package com.example;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;
import java.sql.ResultSet;

import java.time.LocalDate;
import java.time.ZoneId;

public class Statistic{
	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String sql;

	private String message;

	public void createTable(){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "CREATE TABLE IF NOT EXISTS statistics ("+
				"user_id LONG NOT NULL UNIQUE," +
				"current_task INTEGER DEFAULT 0," +
				"current_score INTEGER DEFAULT 0," +
				"baddest_task INTEGER DEFAULT '0'," +
				"baddest_score INTEGER DEFAULT 0," +
				"better_task INTEGER DEFAULT '0'," +
				"better_score INTEGER DEFAULT 0," +
				"last_Active_Date INTEGER  DEFAULT '0'," +
				"streak INTEGER DEFAULT '0'" +
			");";
			Statement stmt = conn.createStatement();
			stmt.execute(sql);
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}
	public String getStatistic(String userName, long userId){
		UserStateManager.getUserState(userId).isActive = false;
		createTable();
		checkStreak(userId);
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT * FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();
			if(result.next()){
				message = "👋Привет " + userName + "!" +
					"	\nТекущее задание: " + result.getInt("current_task") +
					"	\nНаихудшая успеваимость: №" + result.getInt("baddest_task") + ", " + result.getInt("baddest_score") +
					"	\nНаилучшая успеваимость: №" + result.getInt("better_task") + ", " + result.getInt("better_score") +
					"\nТы занимаешься уже " + result.getInt("streak") + " " + getDayForm(result.getInt("streak")) + " подряд!👏";
			}
			else{
				message = "Ваша статистика  пуста.\nВыберите задания чтоб заполнить её!";
			}
		}
		catch(SQLException e){
			e.printStackTrace();
			message = "Таблица статистика пуста.\nВыберите задания чтоб заполнить её!";

		}
		return message;
	}
	public String chooseExercise(long userId, String message){
		try{
			int numberOfExercise = Integer.parseInt(message);
			if(numberOfExercise <= 26 && numberOfExercise >= 1){
				try(Connection conn = DriverManager.getConnection(url)){
					sql = "INSERT INTO statistics (user_id, current_task) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET current_task = excluded.current_task";
					PreparedStatement pstmt = conn.prepareStatement(sql);
					pstmt.setLong(1, userId);
					pstmt.setInt(2, numberOfExercise);
					pstmt.executeUpdate();			
				}
				return "Текущее задание: №" + numberOfExercise;
			}
			else{
				UserStateManager.getUserState(userId).isChoosing = true;
				return "Выберите задание от 1 до 26: ";
			}
		}
		catch(NumberFormatException | SQLException e){
			e.printStackTrace();
			if(!message.isEmpty()) {
				if(e instanceof NumberFormatException){
					UserStateManager.getUserState(userId).isChoosing = true;
					return "Введите номер задания числом: ";
				}
				else{
					return "Произошла ошибка при работе с базой данных";
				}
			}
			return "Перенапровляемся...";
		}
	}
	public String checkStreak(long userId){
		UserState userState = UserStateManager.getUserState(userId);

		String streakMessage = "";
		int currentDate = (int) LocalDate.now(ZoneId.of("Europe/Moscow")).toEpochDay();
		System.out.println(LocalDate.now(ZoneId.of("Europe/Moscow")));
		int lastActiveDate = 0;
		int difference = 0;
		int streak = 0;
		

    try(Connection conn = DriverManager.getConnection(url)){
  		sql = "SELECT last_Active_Date, streak FROM statistics WHERE user_id = ?";
   		PreparedStatement pstmt = conn.prepareStatement(sql);
   		pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();
			
			while(result.next()){
    		lastActiveDate = result.getInt("last_Active_Date");
				streak = result.getInt("streak");
      }

			difference = currentDate - lastActiveDate;
			System.out.println(difference);
			System.out.println(userState.isActive);
			
			if(difference == 1 && userState.isActive || streak == 0 && userState.isActive){
				streak += 1;
				lastActiveDate = currentDate;
				streakMessage = "Стрик обновлён! Вы занимаетесь " + streak + " " + getDayForm(streak);
			}
			else if(difference > 1){
				streak = 0;
			}	

			sql = "UPDATE statistics SET last_Active_Date = ?, streak = ? WHERE user_id = ?";
			pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, lastActiveDate);
			pstmt.setInt(2, streak);
			pstmt.setLong(3, userId);      
			pstmt.executeUpdate();
    }
		catch(SQLException e){
			e.printStackTrace();  
		}
		return streakMessage;
  }
	
	private String getDayForm(int number){
		int lastDigit = number % 10;
		int lastTwoDigits = number % 100;
		if(lastDigit == 1 && lastTwoDigits != 11){
			return "день";
		}
		else if(lastDigit >= 2 && lastDigit <= 4 || !(lastTwoDigits >= 12 && lastTwoDigits <= 14)){
			return "дня";
		}
		else{
			return "дней";
		}
	}
}
