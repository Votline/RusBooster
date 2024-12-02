import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.ResultSet;
import java.util.Random;

import java.util.List;
import java.util.ArrayList;

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
	
	private int task_Id;
	private String sql;

	private static String answer;
	private String returnMessage;
	private Random random = new Random();
	private List<TaskMap> wordsForTask = new ArrayList<>();

	public String makeTask(long userId){
		try(Connection connSet = DriverManager.getConnection(urlStat)){
			sql = "SELECT current_task FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = connSet.prepareStatement(sql);
			pstmt.setLong(1, userId);
			task_Id = pstmt.executeQuery().getInt("current_task");
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		try(Connection conn = DriverManager.getConnection(url)){
			wordsForTask.clear();
			sql = "SELECT word, explanation FROM words WHERE task_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, task_Id);
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
		switch(task_Id){
			case 9:
				Number9 newTask = new Number9();
				returnMessage = newTask.createTask(random, wordsForTask);
				
				break;
			default:
				returnMessage = "Такого задания ещё нет в RusBooster";
				System.out.println(task_Id);
				break;
		}
		return returnMessage;
	}

	public String explanationTask(){
		String message = "Ответ на задание №" + task_Id + ": " + answer;
		return message;
	}

	public static String findAnswer(String values){
		answer = "";
		String[] rows = values.strip().split("\n");
		
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
			if(allMatch) answer += (i+1 + "\n" + row);
		}
		if(answer.isEmpty()) answer = "Совпадений нет. Верный ответ: 0";
		return answer;
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
class Number9{
	public String createTask(Random random, List<TaskMap> task){
		String message = "\n";
		String values = "\n";
		for(int i = 1; i <= 5; i++){
			message += String.valueOf(i) + ") ";
			for(int j = 1; j <= 3; j++){
				int randomIndex = random.nextInt(task.size());
				String key = task.get(randomIndex).key;
				String value = task.get(randomIndex).value;
				task.remove(randomIndex);
				message += key + " ";
				values += value + " ";
			}	
			message += "\n";
			values += "\n";

		}
		System.out.println(message);
		//System.out.println(values);
		System.out.println(MainFunctional.findAnswer(values));
		return message;
	}
}
