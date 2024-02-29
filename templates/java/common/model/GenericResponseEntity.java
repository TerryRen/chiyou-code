package {xxx}.common.model;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.fasterxml.jackson.databind.annotation.JsonNaming;
import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Data;

import java.io.Serializable;
import java.util.Date;
import java.util.UUID;
import java.util.function.Supplier;

/**
 * 响应对象
 *
 * @author mmc
 */
@Data
@JsonIgnoreProperties(ignoreUnknown = true)
@JsonNaming(PropertyNamingStrategies.UpperCamelCaseStrategy.class)
@Schema(description = "泛型风格的响应对象, 实际的数据通过泛型字段data返回")
public class GenericResponseEntity<T> implements Serializable {
    private static final long serialVersionUID = -6547410396516217748L;

    @JsonFormat(shape = JsonFormat.Shape.STRING)
    @Schema(description = "请求时间", example = "2021-08-20T08:13:28.551+0000")
    private final Date requestTime = new Date();

    @Schema(description = "响应Id, 随机生成", example = "52e0bf0d-9022-4726-a6e5-f8686cb28889")
    private final String responseId = UUID.randomUUID().toString();

    @JsonFormat(shape = JsonFormat.Shape.STRING)
    @Schema(description = "响应时间", example = "2021-08-20T08:13:30.551+0000")
    private Date responseTime;

    @Schema(description = "HTTP Status", example = "200")
    private int code = 0;

    @Schema(description = "业务码", example = "A-01-0001")
    private String businessCode = "0-00-0000";

    @Schema(description = "提示消息", example = "success")
    private String msg;

    /**
     * 响应对象
     */
    @Schema(description = "响应数据", example = "{}", type = "Any")
    private T data;

    public GenericResponseEntity() {
        this(200, "success");
    }

    public GenericResponseEntity(int code, String msg) {
        this.code = code;
        this.msg = msg;
        this.responseTime = new Date();
    }

    public T getData() {
        return data;
    }

    public void setData(T data) {
        this.data = data;
        this.responseTime = new Date();
    }

    public void setCode(int code) {
        this.code = code;
        this.responseTime = new Date();
    }

    public static <T> GenericResponseEntity<T> success(T data) {
        GenericResponseEntity<T> response = new GenericResponseEntity<>();
        response.setData(data);
        return response;
    }

    /**
     * success(T data) 方法并不能准确生成requestTime和responseTime
     * 因为在setData时，业务逻辑已经执行完了，并没有计算到真正的业务处理时间
     * 所以此方法就是接受业务逻辑方法，执行完业务方法后才会设置时间
     *
     * @param s   要执行的操作
     * @param <T> 要包装的数据类型
     * @author Chandler.C.Ren
     */
    public static <T> GenericResponseEntity<T> success(Supplier<T> s) {
        GenericResponseEntity<T> response = new GenericResponseEntity<>();
        response.setData(s.get());
        return response;
    }

    //------------------------------------------------以下方法提供链式编程--------------------------------------------------


    public GenericResponseEntity<T> code(int code) {
        setCode(code);
        return this;
    }

    public GenericResponseEntity<T> msg(String msg) {
        setMsg(msg);
        return this;
    }

    public GenericResponseEntity<T> data(T t) {
        setData(t);
        return this;
    }
}

