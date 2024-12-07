import org.telegram.telegrambots.meta.api.methods.send.SendMessage;

public class AdminCommands{
	private Words words = new Words();
	private SendMessage message = new SendMessage();

	public SendMessage dataBase(String messageText, long chatId){
		message.setChatId(String.valueOf(chatId));
		String[] parts = messageText.split(" ", 5);
		if(messageText.contains("/add")){
			String task_id = parts[2];
			String wordName = parts[3];
			String wordExplanation = parts[4];
			message.setText(words.addWord(wordName, wordExplanation, task_id));
		}
		else if(messageText.contains("/remove")){
			String task_id = parts[2];
			String wordName = parts[3];
			message.setText(words.removeWord(wordName, task_id));

		}
		else if(messageText.contains("/showall")){
			String task_id = parts[2];
			message.setText(words.showAllBase(task_id));
		}
		else{
			String text = "Команды: " + "\"" + parts[1] + "\"" +  " не существует";
			message.setText(text);
		}
		return message;
	}
}
