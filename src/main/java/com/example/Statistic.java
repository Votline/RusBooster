package com.example;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;
import java.sql.ResultSet;

import java.time.LocalDate;

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
					"\nТы занимаешься уже " + result.getInt("streak") + " дней подряд!👏";
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
	public void checkStreak(long userId){
		UserState userState = UserStateManager.getUserState(userId);

		int currentDate = (int) LocalDate.now().toEpochDay();
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
			if(difference == 1 && userState.isActive){
				streak += 1;
			}
			else if(difference > 1){
				streak = 0;
			}
			else{
				return;
			}
			lastActiveDate = currentDate;
			

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
  }
}
