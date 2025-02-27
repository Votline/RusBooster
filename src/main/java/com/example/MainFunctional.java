package com.example;

import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.KeyboardRow;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.ReplyKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.User;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.InlineKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.KeyboardButton;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardButton;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.ResultSet;

import java.util.Collections;
import java.util.ArrayList;
import java.util.Random;
import java.util.List;

class TaskMap{
	protected String key;
	protected String value;

	TaskMap(String newKey, String newValue){
		this.key = newKey; this.value = newValue;
	}
}

public class MainFunctional{
	private String urlStat = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/words.db";
	
	private Random random = new Random();
	private Statistic Statistic = new Statistic();
	
	public SendMessage makeTask(long userId){
		UserState userState = UserStateManager.getUserState(userId);
		List<TaskMap> wordsForTask = new ArrayList<>();
		String returnMessage = "";
		String sql = "";
		int taskId = 0;

		try(Connection connSet = DriverManager.getConnection(urlStat)){
			sql = "SELECT current_task FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = connSet.prepareStatement(sql);
			pstmt.setLong(1, userId);
			taskId = pstmt.executeQuery().getInt("current_task");
			userState.currentTask = taskId;
		}
		catch(SQLException e){
			e.printStackTrace();
		}

		try(Connection conn = DriverManager.getConnection(url)){
			wordsForTask.clear();
			sql = "SELECT word, explanation FROM words WHERE task_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, taskId);
			ResultSet result = pstmt.executeQuery();

			while(result.next()){
				String word = result.getString("word");
				String explanation = result.getString("explanation");
				wordsForTask.add(new TaskMap(word, explanation));
			}
		}
		catch(SQLException e){
			e.printStackTrace();
		}

		switch(taskId){
			case 9:
			case 10:
				TaskResult result = Number9to12.createTask(this, random, wordsForTask, 3, userId);
				returnMessage = result.message;
				userState.explanations = result.explanations;
				userState.outputAnswer = result.answer;
				break;
			case 11:
			case 12:
				result = Number9to12.createTask(this, random, wordsForTask, 2, userId);
				returnMessage = result.message;
				userState.explanations = result.explanations;
				userState.outputAnswer = result.answer;
				break;
			default:
				returnMessage = "Такого задания ещё нет в RusBooster";
				break;
		}
		SendMessage sendMessage = new SendMessage();
		sendMessage.setChatId(String.valueOf(userId));
		sendMessage.setText(returnMessage);

		InlineKeyboardMarkup cancelKeyboard = new InlineKeyboardMarkup();
		InlineKeyboardButton cancelButton = new InlineKeyboardButton();
		
		cancelButton.setText("Отказаться от задания");
		cancelButton.setCallbackData("cancelTask");
		
		cancelKeyboard.setKeyboard(Collections.singletonList(Collections.singletonList(cancelButton)));
		sendMessage.setReplyMarkup(cancelKeyboard);

		return sendMessage;
	}

	public SendMessage responseForTask(long chatId, long userId, String message){
		UserState userState = UserStateManager.getUserState(userId);
		userState.isActive = true;

		String additionalMessage = Statistic.checkStreak(userId);
		String sendMessageText = "";

		SendMessage sendMessage = new SendMessage();
		sendMessage.setChatId(String.valueOf(chatId));
		if(additionalMessage != null){	
			sendMessageText = "Ответ на задание №" + userState.currentTask + ": " + userState.outputAnswer + "\n\n" + additionalMessage;
		}
		else {
			sendMessageText = "Ответ на задание №" + userState.currentTask + ": " + userState.outputAnswer;
		}
		sendMessage.setText(sendMessageText);

		InlineKeyboardButton showAllExplanations = new InlineKeyboardButton();
		showAllExplanations.setText("Показать пояснения всех слов");
		showAllExplanations.setCallbackData("showExplanations");

		InlineKeyboardMarkup keyboard = new InlineKeyboardMarkup();
		keyboard.setKeyboard(Collections.singletonList(Collections.singletonList(showAllExplanations)));

		try{
			int number = Integer.parseInt(message);
			int userAnswer = 0;
			while(number != 0){
				userAnswer += number % 10;
				number /= 10;}
			checkAnswer(userAnswer, userId);
			sendMessage.setReplyMarkup(keyboard);
		}
		catch(NumberFormatException e){
			e.printStackTrace();
			sendMessage.setText("Укажите варианты ответов числом");
			UserStateManager.getUserState(userId).isChecking = true;
		}
		Statistic.checkStreak(userId);
		return sendMessage;
	}


	private void checkAnswer(int userAnswer, long userId){
		UserState userState = UserStateManager.getUserState(userId);

		try(Connection conn = DriverManager.getConnection(urlStat)){
			int current_task = 0, current_score = 0, baddest_task = 0, baddest_score = 0, better_task = 0, better_score = 0;
			String sql = "SELECT current_task, current_score, baddest_task, baddest_score, better_task, better_score FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();

			while(result.next()){
				current_task = result.getInt("current_task");
				current_score = result.getInt("current_score");	
				baddest_task = result.getInt("baddest_task");	
				baddest_score = result.getInt("baddest_score");	
				better_task = result.getInt("better_task");	
				better_score = result.getInt("better_score");	
			}

			if(userAnswer == userState.answer){current_score += 1;}
			else{current_score -= 1;}

			if(current_score < baddest_score){
				baddest_score = current_score;
				baddest_task = current_task;
			}
			if(current_score > better_score){
				better_score = current_score;
				better_task = current_task;
			}

			sql = "UPDATE statistics SET current_score = ?, baddest_task = ?, baddest_score = ?, better_task = ?, better_score = ? WHERE user_id = ?";
			PreparedStatement pstmtUpd = conn.prepareStatement(sql);
			pstmtUpd.setInt(1, current_score);
			pstmtUpd.setInt(2, baddest_task);
			pstmtUpd.setInt(3, baddest_score);
			pstmtUpd.setInt(4, better_task);
			pstmtUpd.setInt(5, better_score);
			pstmtUpd.setLong(6, userId);
			pstmtUpd.executeUpdate();
		}
		catch(SQLException e ){
			e.printStackTrace();
		}
	}

	public String findAnswer(String values, long userId){
		UserState userState = UserStateManager.getUserState(userId);
		String output = "";
		
		if(values == null){
			output = "Ошибка, нет значений для поиска";
			return output;
		}

		userState.answer = 0;
		String answerForOutput = "";
		String[] rows = values.strip().split("\n\n\n");
		
		for(int i = 0; i < rows.length; i++){
			String row = rows[i];
			char base = findUpperCase(getBeforeDot(row));
			if(base == '\0') continue;
			boolean allMatch = true;
			
			while(!row.isEmpty()){
				String currentWord = getBeforeDot(row);
				if(base != findUpperCase(currentWord)){
					allMatch = false;
					break;
				}
				row = removeBeforeDot(row);
			}
			if(allMatch) {
				output += "\n" + rows[i];
				userState.answer += i+1;
				answerForOutput += String.valueOf(i+1);
			}
		}
		if(output.isEmpty()) {
			output = "Совпадений нет. Верный ответ: 0";
			userState.answer = 0;
		}
		else{
			output = answerForOutput + output;
		}

		return output;
	}


	
	public SendMessage showExplanations(long chatId){
		String savedExplanations = UserStateManager.getUserState(chatId).explanations;

		SendMessage sendMessage = new SendMessage();
		sendMessage.setChatId(String.valueOf(chatId));
		sendMessage.setText(savedExplanations != null ? savedExplanations : "Пояснения отсутствуют");

		ReplyKeyboardMarkup keyboardMarkup = new ReplyKeyboardMarkup();
		keyboardMarkup.setResizeKeyboard(true);

		List<KeyboardRow> keyboard = new ArrayList<>();
		KeyboardRow row = new KeyboardRow();
		row.add(new KeyboardButton("Выбрать задание"));
		row.add(new KeyboardButton("Проверить знания"));
		row.add(new KeyboardButton("Статистика"));
		keyboard.add(row);
		
		keyboardMarkup.setKeyboard(keyboard);
		sendMessage.setReplyMarkup(keyboardMarkup);

		return sendMessage;
	}

	private static char findUpperCase(String text){
		for(char c : text.toCharArray()){
			if(Character.isUpperCase(c)) return c;
		}
		return '\0';
	}
	private static String getBeforeDot(String text){
		return text.substring(0, text.indexOf(".")).strip().split(" ")[0];
	}
	private static String removeBeforeDot(String text){
		return text.substring(text.indexOf(".") + 1).strip();
	}



}

class TaskResult {
	String message;
	String answer;
	String explanations;
	
	TaskResult(String message, String answer, String explanations) {
		this.message = message;
		this.answer = answer;
		this.explanations = explanations;
	}
}

class Number9to12 {
	public static TaskResult createTask(MainFunctional functional, Random random, List<TaskMap> task, int wordCount, long userId) {
		String message = "Укажите варианты ответов, в которых во всех словах одного ряда пропущена одна и та же буква. Запишите номера ответов.\n";
		String explanations = "";
		
		for(int i = 1; i <= 5; i++) {
			message += String.valueOf(i) + ") ";
			for(int j = 1; j <= wordCount; j++) {
				int randomIndex = random.nextInt(task.size());
				String key = task.get(randomIndex).key;
				String value = task.get(randomIndex).value;
				task.remove(randomIndex);
				message += key + " ";
				explanations += "\n" + value + "\n";
			}	
			message += "\n";
			explanations += "\n";
		}
				
		String answer = functional.findAnswer(explanations, userId);
		if(answer.equals("Совпадений нет. Верный ответ: 0")){
			return createTask(functional, random, task, wordCount, userId);
		}
		return new TaskResult(message, answer, explanations);
	}
}
