import org.telegram.telegrambots.meta.api.methods.send.SendMessage;

public class AdminCommands{
	private Words words = new Words();
	private SendMessage message = new SendMessage();

	public SendMessage dataBase(String messageText, long chatId){
		message.setChatId(String.valueOf(chatId));
		String[] parts = messageText.split(" ", 4);
		if(messageText.contains("/add")){
			String wordName = parts[2];
			String wordExplanation = parts[3];
			message.setText(words.addWord(wordName, wordExplanation));
		}
/*		else if(messageText.contains("/remove")){
			String wordName = parts[2];
			message.setText(words.removeWord(wordName));

		}
		else if(messageText.contains("/showall")){
			message.setText(words.showAllBase());
		}
		else{
			String text = "Команды: " + "\"" + parts[1] + "\"" +  " не существует";
			System.out.println("Вывел ошибку: " + text);
			message.setText(text);
		}
*/		return message;
	}
}
